package model

import (
	"github.com/figment-networks/celo-indexer/types"
	"time"
)

const (
	SyncableStatusRunning SyncableStatus = iota + 1
	SyncableStatusCompleted
)

type SyncableStatus int

type Syncable struct {
	*ModelWithTimestamps

	Height      int64       `json:"height"`
	Time        *types.Time `json:"time"`
	ChainId     uint64      `json:"chain_id"`
	Epoch       *int64      `json:"epoch"`
	LastInEpoch *bool       `json:"last_in_epoch"`

	IndexVersion  int64          `json:"index_version"`
	Status        SyncableStatus `json:"status"`
	ReportID      types.ID       `json:"report_id"`
	StartedAt     types.Time     `json:"started_at"`
	ProcessedAt   *types.Time    `json:"processed_at"`
	Duration      time.Duration  `json:"duration"`
	RequestsCount uint64         `json:"requests_count"`
}

// - Methods
func (Syncable) TableName() string {
	return "syncables"
}

func (s *Syncable) Valid() bool {
	return s.Height >= 0
}

func (s *Syncable) Equal(m Syncable) bool {
	return s.Height == m.Height
}

func (s *Syncable) SetStatus(newStatus SyncableStatus) {
	s.Status = newStatus
}

func (s *Syncable) MarkProcessed(indexVersion int64, requestsCount uint64) {
	t := types.NewTimeFromTime(time.Now())
	duration := time.Since(s.StartedAt.Time)

	s.Status = SyncableStatusCompleted
	s.Duration = duration
	s.ProcessedAt = t
	s.IndexVersion = indexVersion
	s.RequestsCount = requestsCount
}
