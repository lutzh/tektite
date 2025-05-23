package control

import (
	"github.com/pkg/errors"
	"github.com/spirit-labs/tektite/common"
	log "github.com/spirit-labs/tektite/logger"
	"github.com/spirit-labs/tektite/lsm"
	"github.com/spirit-labs/tektite/objstore"
	"github.com/spirit-labs/tektite/sst"
	"time"

	"sync"
	"sync/atomic"
)

/*
LsmHolder is a wrapper around the Lsm manager which handles queueing of apply changes requests and persistence to
object storage
*/
type LsmHolder struct {
	lock                   sync.RWMutex
	objStore               objstore.Client
	lsmOpts                lsm.Conf
	metaDataBucketName     string
	metaDataKey            string
	lsmManager             *lsm.Manager
	started                bool
	stopping               atomic.Bool
	hasQueuedRegistrations atomic.Bool
	queuedRegistrations    []queuedRegistration
	waitingCompletions     []func(error) error
	stateWriteTimer        *time.Timer
	stateWriteInterval     time.Duration
	metaDataEtag           string
}

type queuedRegistration struct {
	regBatch       lsm.RegistrationBatch
	completionFunc func(error) error
}

func NewLsmHolder(metaDataBucketName string,
	metaDataKey string, objStoreClient objstore.Client, stateWriteInterval time.Duration, lsmOpts lsm.Conf) *LsmHolder {
	return &LsmHolder{
		objStore:           objStoreClient,
		lsmOpts:            lsmOpts,
		metaDataBucketName: metaDataBucketName,
		metaDataKey:        metaDataKey,
		stateWriteInterval: stateWriteInterval,
	}
}

func (s *LsmHolder) Start() error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.started {
		return nil
	}
	metaData, metaDataEtag, err := s.loadMetadata()
	if err != nil {
		return err
	}
	lsmManager := lsm.NewManager(s.objStore, s.maybeRetryApplies, true, false, s.lsmOpts)
	if err := lsmManager.Start(metaData); err != nil {
		return err
	}
	s.lsmManager = lsmManager
	s.metaDataEtag = metaDataEtag
	s.scheduleStateWriteTimer()
	s.started = true
	return nil
}

const objectStoreCallTimeout = 5 * time.Second

func (s *LsmHolder) loadMetadata() ([]byte, string, error) {
	for !s.stopping.Load() {
		objectInfo, exists, err := objstore.GetObjectInfoWithTimeout(s.objStore, s.metaDataBucketName, s.metaDataKey,
			objectStoreCallTimeout)
		if err == nil {
			if !exists {
				return nil, "", nil
			}
			metaData, err := objstore.GetWithTimeout(s.objStore, s.metaDataBucketName, s.metaDataKey, objectStoreCallTimeout)
			if err == nil {
				return metaData, objectInfo.Etag, nil
			}
		}
		if common.IsUnavailableError(err) {
			log.Debugf("object store is unavailable on load metadata etag, will retry: %v", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		return nil, "", err
	}
	return nil, "", errors.New("lsmHolder is stopping")
}

func (s *LsmHolder) Stop() error {
	s.stopping.Store(true)
	s.lock.Lock()
	defer s.lock.Unlock()
	if !s.started {
		return nil
	}
	return s.stop()
}

func (s *LsmHolder) stop() error {
	if err := s.lsmManager.Stop(); err != nil {
		return err
	}
	if s.stateWriteTimer != nil {
		s.stateWriteTimer.Stop()
	}
	s.started = false
	return nil
}

func (s *LsmHolder) scheduleStateWriteTimer() {
	s.stateWriteTimer = time.AfterFunc(s.stateWriteInterval, func() {
		s.lock.Lock()
		unlocked := false
		defer func() {
			if !unlocked {
				s.lock.Unlock()
			}
		}()
		if !s.started {
			return
		}
		completions, err := s.maybeWriteState()
		s.lock.Unlock()
		unlocked = true
		// Call completions - these must be called outside the lock to prevent deadlock with offsets cache locking
		for _, cf := range completions {
			if err2 := cf(err); err2 != nil {
				log.Errorf("failed to apply completion function: %v", err2)
			}
		}
		s.lock.Lock()
		defer s.lock.Unlock()
		s.scheduleStateWriteTimer()
	})
}

func (s *LsmHolder) GetTablesForHighestKeyWithPrefix(prefix []byte) ([]sst.SSTableID, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if err := s.checkStarted(); err != nil {
		return nil, err
	}
	return s.lsmManager.GetTablesForHighestKeyWithPrefix(prefix)
}

// ApplyLsmChanges - apply some changes to the LSM structure. Note that this method completes asynchronously as
// L0 registrations will be queued if there is not enough free space
func (s *LsmHolder) ApplyLsmChanges(regBatch lsm.RegistrationBatch, completionFunc func(error) error) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if err := s.checkStarted(); err != nil {
		return completionFunc(err)
	}
	ok, err := s.lsmManager.ApplyChanges(regBatch, false)
	if err != nil {
		return completionFunc(err)
	}
	if ok {
		s.waitingCompletions = append(s.waitingCompletions, completionFunc)
		return nil
	}
	// L0 is full - queue the registration - it will be retried when there is space in L0
	s.hasQueuedRegistrations.Store(true)
	s.queuedRegistrations = append(s.queuedRegistrations, queuedRegistration{
		regBatch:       regBatch,
		completionFunc: completionFunc,
	})
	return nil
}

func (s *LsmHolder) maybeRetryApplies() {
	// check atomic outside lock to reduce contention
	if !s.hasQueuedRegistrations.Load() {
		return
	}
	if err := s.maybeRetryApplies0(); err != nil {
		log.Errorf("failed to retry applies: %v", err)
	}
}

func (s *LsmHolder) maybeWriteState() ([]func(error) error, error) {
	if len(s.waitingCompletions) == 0 {
		return nil, nil
	}
	metaData := s.lsmManager.GetMasterRecordBytes()
	var ok bool
	var err error
	var etag string
	if s.metaDataEtag == "" {
		// First time - no state yet
		ok, etag, err = objstore.PutIfNotExistsWithTimeout(s.objStore, s.metaDataBucketName, s.metaDataKey, metaData,
			objectStoreCallTimeout)
	} else {
		// Put only succeeds if etag hasn't changed
		ok, etag, err = objstore.PutIfMatchingEtagWithTimeout(s.objStore, s.metaDataBucketName, s.metaDataKey, metaData,
			s.metaDataEtag, objectStoreCallTimeout)
	}
	if err == nil && !ok {
		// Failed to store data as another controller has saved state and changed the etag
		err = common.NewTektiteErrorf(common.Unavailable, "controller not leader")
		// No longer leader so stop the controller
		if err := s.stop(); err != nil {
			log.Warnf("failed to stop controller: %v", err)
		}
	}
	if err != nil {
		log.Warnf("failed to store lsm state: %v", err)
	} else {
		s.metaDataEtag = etag
	}
	res := make([]func(error) error, len(s.waitingCompletions))
	copy(res, s.waitingCompletions)
	s.waitingCompletions = nil
	return res, err
}

func (s *LsmHolder) maybeRetryApplies0() error {
	s.lock.Lock()
	defer s.lock.Unlock()
	var newWaiting []queuedRegistration
	// Try and reapply the registrations
	for _, queuedReg := range s.queuedRegistrations {
		ok, err := s.lsmManager.ApplyChanges(queuedReg.regBatch, false)
		if err != nil {
			return queuedReg.completionFunc(err)
		}
		if ok {
			// Added to L0 ok. Now add to the waiting completions - these will be called once the metadata is committed
			// to object storage
			s.waitingCompletions = append(s.waitingCompletions, queuedReg.completionFunc)
		} else {
			// L0 is full again - keep the queued registration
			newWaiting = append(newWaiting, queuedReg)
		}
	}
	s.queuedRegistrations = newWaiting
	return nil
}

func (s *LsmHolder) QueryTablesInRange(keyStart []byte, keyEnd []byte) (lsm.OverlappingTables, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if err := s.checkStarted(); err != nil {
		return nil, err
	}
	return s.lsmManager.QueryTablesInRange(keyStart, keyEnd)
}

func (s *LsmHolder) checkStarted() error {
	if !s.started {
		return common.NewTektiteErrorf(common.Unavailable, "lsm holder is not started")
	}
	return nil
}
