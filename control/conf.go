package control

import (
	"github.com/spirit-labs/tektite/common"
	"github.com/spirit-labs/tektite/lsm"
	"time"
)

type Conf struct {
	ControllerMetaDataBucketName string
	ControllerMetaDataKey        string
	SSTableBucketName            string
	DataFormat                   common.DataFormat
	TableNotificationInterval    time.Duration
	LsmConf                      lsm.Conf
	SequencesBlockSize           int
	AzInfo                       string
	LsmStateWriteInterval        time.Duration
}

func NewConf() Conf {
	return Conf{
		ControllerMetaDataBucketName: "controller-meta-data",
		ControllerMetaDataKey:        "controller-meta-data",
		SSTableBucketName:            "tektite-data",
		DataFormat:                   common.DataFormatV1,
		TableNotificationInterval:    5 * time.Second,
		LsmConf:                      lsm.NewConf(),
		SequencesBlockSize:           100,
		LsmStateWriteInterval:        10 * time.Millisecond,
	}
}

func (c *Conf) Validate() error {
	if err := c.LsmConf.Validate(); err != nil {
		return err
	}
	return nil
}
