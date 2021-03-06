package indexer

import (
	"context"
	"fmt"
	"github.com/figment-networks/celo-indexer/store"
	"github.com/figment-networks/celo-indexer/store/psql"

	"github.com/figment-networks/celo-indexer/client"
	"github.com/figment-networks/celo-indexer/config"
	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/pkg/errors"
)

var (
	_ pipeline.Source = (*backfillSource)(nil)
)

func NewBackfillSource(cfg *config.Config, client client.Client, syncableDb store.Syncables, indexVersion int64) (*backfillSource, error) {
	src := &backfillSource{
		cfg:    cfg,
		client: client,

		syncableDb: syncableDb,

		currentIndexVersion: indexVersion,
	}

	if err := src.init(); err != nil {
		return nil, err
	}

	return src, nil
}

type backfillSource struct {
	cfg    *config.Config
	client client.Client

	syncableDb store.Syncables

	currentIndexVersion int64

	currentHeight int64
	startHeight   int64
	endHeight     int64
	err           error
}

func (s *backfillSource) Skip(pipeline.StageName) bool {
	return false
}

func (s *backfillSource) Next(context.Context, pipeline.Payload) bool {
	if s.err == nil && s.currentHeight < s.endHeight {
		s.currentHeight = s.currentHeight + 1
		return true
	}
	return false
}

func (s *backfillSource) Current() int64 {
	return s.currentHeight
}

func (s *backfillSource) Err() error {
	return s.err
}

func (s *backfillSource) Len() int64 {
	return s.endHeight - s.startHeight + 1
}

func (s *backfillSource) init() error {
	if err := s.setStartHeight(); err != nil {
		return err
	}
	if err := s.setEndHeight(); err != nil {
		return err
	}
	return nil
}

func (s *backfillSource) setStartHeight() error {
	syncable, err := s.syncableDb.FindFirstByDifferentIndexVersion(s.currentIndexVersion)
	if err != nil {
		if err == psql.ErrNotFound {
			return errors.New(fmt.Sprintf("nothing to backfill [currentIndexVersion=%d]", s.currentIndexVersion))
		}
		return err
	}

	s.currentHeight = syncable.Height
	s.startHeight = syncable.Height
	return nil
}

func (s *backfillSource) setEndHeight() error {
	syncable, err := s.syncableDb.FindMostRecentByDifferentIndexVersion(s.currentIndexVersion)
	if err != nil {
		if err == psql.ErrNotFound {
			return errors.New(fmt.Sprintf("nothing to backfill [currentIndexVersion=%d]", s.currentIndexVersion))
		}
		return err
	}

	s.endHeight = syncable.Height
	return nil
}
