package indexing

import (
	"context"
	"fmt"
	"time"

	"github.com/figment-networks/celo-indexer/metrics"
	"github.com/figment-networks/celo-indexer/store/psql"

	"github.com/figment-networks/celo-indexer/config"
	"github.com/figment-networks/celo-indexer/indexer"
	"github.com/figment-networks/celo-indexer/types"
	"github.com/figment-networks/celo-indexer/utils/logger"
	"github.com/pkg/errors"
)

var (
	ErrPurgingDisabled = errors.New("purging disabled")
)

type purgeUseCase struct {
	cfg *config.Config
	db  *psql.Store
}

func NewPurgeUseCase(cfg *config.Config, db *psql.Store) *purgeUseCase {
	return &purgeUseCase{
		cfg: cfg,
		db:  db,
	}
}

func (uc *purgeUseCase) Execute(ctx context.Context) error {
	defer metrics.LogUsecaseDuration(time.Now(), "purge")

	configParser, err := indexer.NewConfigParser(uc.cfg.IndexerConfigFile)
	if err != nil {
		return err
	}
	currentIndexVersion := configParser.GetCurrentVersionId()

	if err := uc.purgeBlocks(currentIndexVersion); err != nil {
		return err
	}

	if err := uc.purgeValidators(currentIndexVersion); err != nil {
		return err
	}

	if err := uc.purgeValidatorGroups(currentIndexVersion); err != nil {
		return err
	}

	return nil
}

func (uc *purgeUseCase) purgeBlocks(currentIndexVersion int64) error {
	if err := uc.purgeBlockSequences(currentIndexVersion); uc.checkErr(err) {
		return err
	}
	if err := uc.purgeBlockSummaries(types.IntervalHourly, uc.cfg.PurgeHourlySummariesInterval); uc.checkErr(err) {
		return err
	}
	return nil
}

func (uc *purgeUseCase) purgeValidators(currentIndexVersion int64) error {
	if err := uc.purgeValidatorSequences(currentIndexVersion); uc.checkErr(err) {
		return err
	}

	if err := uc.purgeValidatorSummaries(types.IntervalHourly, uc.cfg.PurgeHourlySummariesInterval); uc.checkErr(err) {
		return err
	}

	return nil
}

func (uc *purgeUseCase) purgeValidatorGroups(currentIndexVersion int64) error {
	if err := uc.purgeValidatorGroupSequences(currentIndexVersion); uc.checkErr(err) {
		return err
	}

	if err := uc.purgeValidatorGroupSummaries(types.IntervalHourly, uc.cfg.PurgeHourlySummariesInterval); uc.checkErr(err) {
		return err
	}

	return nil
}

func (uc *purgeUseCase) purgeBlockSequences(currentIndexVersion int64) error {
	blockSeq, err := uc.db.GetBlocks().BlockSeq.FindMostRecent()
	if err != nil {
		return err
	}
	lastSeqTime := blockSeq.Time.Time

	duration, err := uc.parseDuration(uc.cfg.PurgeSequencesInterval)
	if err != nil {
		if err == ErrPurgingDisabled {
			logger.Info("purging block sequences disabled. Purge interval set to 0.")
		}
		return err
	}

	purgeThresholdFromLastSeq := lastSeqTime.Add(-*duration)

	activityPeriods, err := uc.db.GetBlocks().BlockSummary.FindActivityPeriods(types.IntervalDaily, currentIndexVersion)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("purging summarized block sequences... [older than=%s]", purgeThresholdFromLastSeq))

	deletedCount, err := uc.db.GetBlocks().BlockSeq.DeleteOlderThan(purgeThresholdFromLastSeq, activityPeriods)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("%d block sequences purged", *deletedCount))

	return nil
}

func (uc *purgeUseCase) purgeBlockSummaries(interval types.SummaryInterval, purgeInterval string) error {
	blockSummary, err := uc.db.GetBlocks().BlockSummary.FindMostRecentByInterval(interval)
	if err != nil {
		return err
	}
	lastSummaryTimeBucket := blockSummary.TimeBucket.Time

	duration, err := uc.parseDuration(purgeInterval)
	if err != nil {
		if err == ErrPurgingDisabled {
			logger.Info(fmt.Sprintf("purging block summaries disabled [interval=%s] [purge_interval=%s]", interval, purgeInterval))
		}
		return err
	}

	purgeThreshold := lastSummaryTimeBucket.Add(-*duration)

	logger.Info(fmt.Sprintf("purging block summaries... [interval=%s] [older than=%s]", interval, purgeThreshold))

	deletedCount, err := uc.db.GetBlocks().BlockSummary.DeleteOlderThan(interval, purgeThreshold)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("%d block summaries purged [interval=%s]", *deletedCount, interval))

	return nil
}

func (uc *purgeUseCase) purgeValidatorSequences(currentIndexVersion int64) error {
	validatorSeq, err := uc.db.GetValidators().ValidatorSeq.FindMostRecent()
	if err != nil {
		return err
	}
	lastSeqTime := validatorSeq.Time.Time

	validatorSummary, err := uc.db.GetValidators().ValidatorSummary.FindMostRecent()
	if err != nil {
		return err
	}
	lastSummaryTimeBucket := validatorSummary.TimeBucket.Time

	duration, err := uc.parseDuration(uc.cfg.PurgeSequencesInterval)
	if err != nil {
		if err == ErrPurgingDisabled {
			logger.Info("purging validator sequences disabled. Purge interval set to 0.")
		}
		return err
	}

	purgeThresholdFromConfig := lastSeqTime.Add(-*duration)

	var purgeThreshold time.Time
	if purgeThresholdFromConfig.Before(lastSummaryTimeBucket) {
		purgeThreshold = purgeThresholdFromConfig
	} else {
		purgeThreshold = lastSummaryTimeBucket
	}

	logger.Info(fmt.Sprintf("purging validator sequences... [older than=%s]", purgeThreshold))

	deletedCount, err := uc.db.GetValidators().ValidatorSeq.DeleteOlderThan(purgeThreshold)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("%d validator sequences purged", *deletedCount))

	return nil
}

func (uc *purgeUseCase) purgeValidatorSummaries(interval types.SummaryInterval, purgeInterval string) error {
	validatorSummary, err := uc.db.GetValidators().ValidatorSummary.FindMostRecentByInterval(interval)
	if err != nil {
		return err
	}
	lastSummaryTimeBucket := validatorSummary.TimeBucket.Time

	duration, err := uc.parseDuration(purgeInterval)
	if err != nil {
		if err == ErrPurgingDisabled {
			logger.Info(fmt.Sprintf("purging validator summaries disabled [interval=%s] [purge_interval=%s]", interval, purgeInterval))
		}
		return err
	}

	purgeThreshold := lastSummaryTimeBucket.Add(-*duration)

	logger.Info(fmt.Sprintf("purging validator summaries... [interval=%s] [older than=%s]", interval, purgeThreshold))

	deletedCount, err := uc.db.GetValidators().ValidatorSummary.DeleteOlderThan(interval, purgeThreshold)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("%d validator summaries purged [interval=%s]", *deletedCount, interval))

	return nil
}

func (uc *purgeUseCase) purgeValidatorGroupSequences(currentIndexVersion int64) error {
	validatorGroupSeq, err := uc.db.GetValidatorGroups().ValidatorGroupSeq.FindMostRecent()
	if err != nil {
		return err
	}
	lastSeqTime := validatorGroupSeq.Time.Time

	validatorGroupSummary, err := uc.db.GetValidatorGroups().ValidatorGroupSummary.FindMostRecent()
	if err != nil {
		return err
	}
	lastSummaryTimeBucket := validatorGroupSummary.TimeBucket.Time

	duration, err := uc.parseDuration(uc.cfg.PurgeSequencesInterval)
	if err != nil {
		if err == ErrPurgingDisabled {
			logger.Info("purging validator group sequences disabled. Purge interval set to 0.")
		}
		return err
	}

	purgeThresholdFromConfig := lastSeqTime.Add(-*duration)

	var purgeThreshold time.Time
	if purgeThresholdFromConfig.Before(lastSummaryTimeBucket) {
		purgeThreshold = purgeThresholdFromConfig
	} else {
		purgeThreshold = lastSummaryTimeBucket
	}

	logger.Info(fmt.Sprintf("purging validator group sequences... [older than=%s]", purgeThreshold))

	deletedCount, err := uc.db.GetValidatorGroups().ValidatorGroupSeq.DeleteOlderThan(purgeThreshold)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("%d validator group sequences purged", *deletedCount))

	return nil
}

func (uc *purgeUseCase) purgeValidatorGroupSummaries(interval types.SummaryInterval, purgeInterval string) error {
	validatorGroupSummary, err := uc.db.GetValidatorGroups().ValidatorGroupSummary.FindMostRecentByInterval(interval)
	if err != nil {
		return err
	}
	lastSummaryTimeBucket := validatorGroupSummary.TimeBucket.Time

	duration, err := uc.parseDuration(purgeInterval)
	if err != nil {
		if err == ErrPurgingDisabled {
			logger.Info(fmt.Sprintf("purging validator group summaries disabled [interval=%s] [purge_interval=%s]", interval, purgeInterval))
		}
		return err
	}

	purgeThreshold := lastSummaryTimeBucket.Add(-*duration)

	logger.Info(fmt.Sprintf("purging validator group summaries... [interval=%s] [older than=%s]", interval, purgeThreshold))

	deletedCount, err := uc.db.GetValidatorGroups().ValidatorGroupSummary.DeleteOlderThan(interval, purgeThreshold)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("%d validator group summaries purged [interval=%s]", *deletedCount, interval))

	return nil
}

func (uc *purgeUseCase) parseDuration(interval string) (*time.Duration, error) {
	duration, err := time.ParseDuration(interval)
	if err != nil {
		return nil, err
	}

	if duration == 0 {
		return nil, ErrPurgingDisabled
	}
	return &duration, nil
}

func (uc *purgeUseCase) checkErr(err error) bool {
	return err != nil && err != ErrPurgingDisabled && err != psql.ErrNotFound
}
