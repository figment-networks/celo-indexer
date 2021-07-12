package store

import (
	"time"

	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/types"
)

type Database interface {
	GetTotalSize() (*GetTotalSizeResult, error)
}

// GetAvgTimesForIntervalRow Contains row of data for FindSummary query
type GetTotalSizeResult struct {
	Size float64 `json:"size"`
}

type Reports interface {
	Create(report *model.Report) error
	Save(report *model.Report) error
	FindNotCompletedByIndexVersion(indexVersion int64, kinds ...model.ReportKind) (*model.Report, error)
	FindNotCompletedByKind(kinds ...model.ReportKind) (*model.Report, error)
	Last() (*model.Report, error)
	DeleteByKinds(kinds []model.ReportKind) error
}

type Syncables interface {
	Save(syncable *model.Syncable) error
	FindSmallestIndexVersion() (*int64, error)
	FindByHeight(height int64) (syncable *model.Syncable, err error)
	FindMostRecent() (*model.Syncable, error)
	FindLastInEpochForHeight(height int64) (syncable *model.Syncable, err error)
	FindLastInEpoch(epoch int64) (syncable *model.Syncable, err error)
	FindFirstByDifferentIndexVersion(indexVersion int64) (*model.Syncable, error)
	FindMostRecentByDifferentIndexVersion(indexVersion int64) (*model.Syncable, error)
	CreateOrUpdate(val *model.Syncable) error
	SetProcessedAtForRange(reportID types.ID, startHeight int64, endHeight int64) error
}

type SystemEvents interface {
	BulkUpsert(records []model.SystemEvent) error
	FindByHeight(height int64) ([]model.SystemEvent, error)
	FindByActor(actorAddress string, query FindSystemEventByActorQuery) ([]model.SystemEvent, error)
	FindAll(query FindAll) ([]model.SystemEvent, error)
	FindUnique(height int64, address string, kind model.SystemEventKind) (*model.SystemEvent, error)
	FindMostRecent() (*model.SystemEvent, error)
	DeleteOlderThan(purgeThreshold time.Time) (*int64, error)
}

type FindSystemEventByActorQuery struct {
	Kind      *model.SystemEventKind
	MinHeight *int64
}

type FindAll struct {
	Page  int64
	Limit int64
}

type Jobs interface {
	Create(jobs []model.Job) error
	Update(job *model.Job, fields ...string) error
	FindByHeight(height int64) (*model.Job, error)
	FindAllUnfinished() ([]model.Job, error)
	LastFinishedHeight() (int64, error)
	LastSyncedHeight() (int64, error)
}
