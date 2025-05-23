package pusher

import (
	"github.com/spirit-labs/tektite/common"
	"github.com/spirit-labs/tektite/compress"
	"time"
)

type Conf struct {
	WriteTimeout                             time.Duration
	AvailabilityRetryInterval                time.Duration
	BufferMaxSizeBytes                       int
	DataFormat                               common.DataFormat
	DataBucketName                           string
	OffsetSnapshotInterval                   time.Duration
	TopicCompactionInterval                  time.Duration
	CompactedTopicLastOffsetSnapshotInterval time.Duration
	EnforceProduceOnLeader                   bool
	TableCompressionType                     compress.CompressionType
}

func NewConf() Conf {
	return Conf{
		BufferMaxSizeBytes:                       DefaultBufferSizeMaxBytes,
		WriteTimeout:                             DefaultWriteTimeout,
		AvailabilityRetryInterval:                DefaultAvailabilityRetryInterval,
		DataFormat:                               DefaultDataFormat,
		DataBucketName:                           DefaultDataBucketName,
		OffsetSnapshotInterval:                   DefaultOffsetSnapshotInterval,
		CompactedTopicLastOffsetSnapshotInterval: DefaultCompactedTopicLastOffsetSnapshotInterval,
	}
}

func (c *Conf) Validate() error {
	return nil
}

const (
	DefaultWriteTimeout                             = 200 * time.Millisecond
	DefaultAvailabilityRetryInterval                = 1 * time.Second
	DefaultBufferSizeMaxBytes                       = 4 * 1024 * 1024
	DefaultDataFormat                               = common.DataFormatV1
	DefaultDataBucketName                           = "tektite-data"
	DefaultOffsetSnapshotInterval                   = 5 * time.Second
	DefaultCompactedTopicLastOffsetSnapshotInterval = 5 * time.Second
)
