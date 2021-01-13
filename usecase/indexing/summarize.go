package indexing

import (
	"context"
	"fmt"
	"github.com/figment-networks/celo-indexer/store/psql"
	"time"

	"github.com/figment-networks/celo-indexer/config"
	"github.com/figment-networks/celo-indexer/indexer"
	"github.com/figment-networks/celo-indexer/metric"
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/types"
	"github.com/figment-networks/celo-indexer/utils/logger"
)

type summarizeUseCase struct {
	cfg *config.Config
	db  *psql.Store
}

func NewSummarizeUseCase(cfg *config.Config, db *psql.Store) *summarizeUseCase {
	return &summarizeUseCase{
		cfg: cfg,
		db:  db,
	}
}

func (uc *summarizeUseCase) Execute(ctx context.Context) error {
	defer metric.LogUseCaseDuration(time.Now(), "summarize")

	configParser, err := indexer.NewConfigParser(uc.cfg.IndexerConfigFile)
	if err != nil {
		return err
	}
	currentIndexVersion := configParser.GetCurrentVersionId()

	if err := uc.summarizeBlockSeq(types.IntervalHourly, currentIndexVersion); err != nil {
		return err
	}

	if err := uc.summarizeBlockSeq(types.IntervalDaily, currentIndexVersion); err != nil {
		return err
	}

	if err := uc.summarizeValidatorSeq(types.IntervalHourly, currentIndexVersion); err != nil {
		return err
	}

	if err := uc.summarizeValidatorSeq(types.IntervalDaily, currentIndexVersion); err != nil {
		return err
	}

	if err := uc.summarizeValidatorGroupSeq(types.IntervalHourly, currentIndexVersion); err != nil {
		return err
	}

	if err := uc.summarizeValidatorGroupSeq(types.IntervalDaily, currentIndexVersion); err != nil {
		return err
	}

	return nil
}

func (uc *summarizeUseCase) summarizeBlockSeq(interval types.SummaryInterval, currentIndexVersion int64) error {
	logger.Info(fmt.Sprintf("summarizing block sequences... [interval=%s]", interval))

	activityPeriods, err := uc.db.GetBlocks().BlockSummary.FindActivityPeriods(interval, currentIndexVersion)
	if err != nil {
		return err
	}

	rawSummaryItems, err := uc.db.GetBlocks().BlockSeq.Summarize(interval, activityPeriods)
	if err != nil {
		return err
	}

	var newModels []model.BlockSummary
	var existingModels []model.BlockSummary
	for _, rawSummary := range rawSummaryItems {
		summary := &model.Summary{
			TimeInterval: interval,
			TimeBucket:   rawSummary.TimeBucket,
			IndexVersion: currentIndexVersion,
		}
		query := model.BlockSummary{
			Summary: summary,
		}

		existingBlockSummary, err := uc.db.GetBlocks().BlockSummary.Find(&query)
		if err != nil {
			if err == psql.ErrNotFound {
				blockSummary := model.BlockSummary{
					Summary: summary,

					Count:        rawSummary.Count,
					BlockTimeAvg: rawSummary.BlockTimeAvg,
				}
				if err := uc.db.GetBlocks().BlockSummary.Create(&blockSummary); err != nil {
					return err
				}
				newModels = append(newModels, blockSummary)
			} else {
				return err
			}
		} else {
			existingBlockSummary.Count = rawSummary.Count
			existingBlockSummary.BlockTimeAvg = rawSummary.BlockTimeAvg

			if err := uc.db.GetBlocks().BlockSummary.Save(existingBlockSummary); err != nil {
				return err
			}
			existingModels = append(existingModels, *existingBlockSummary)
		}
	}

	logger.Info(fmt.Sprintf("block sequences summarized [created=%d] [updated=%d]", len(newModels), len(existingModels)))

	return nil
}

func (uc *summarizeUseCase) summarizeValidatorSeq(interval types.SummaryInterval, currentIndexVersion int64) error {
	logger.Info(fmt.Sprintf("summarizing validator sequences... [interval=%s]", interval))

	activityPeriods, err := uc.db.GetValidators().ValidatorSummary.FindActivityPeriods(interval, currentIndexVersion)
	if err != nil {
		return err
	}

	rawSeqSummaryItems, err := uc.db.GetValidators().ValidatorSeq.Summarize(interval, activityPeriods)
	if err != nil {
		return err
	}

	var newModels []model.ValidatorSummary
	var existingModels []model.ValidatorSummary
	for _, rawSeqSummaryItem := range rawSeqSummaryItems {
		summary := &model.Summary{
			TimeInterval: interval,
			TimeBucket:   rawSeqSummaryItem.TimeBucket,
			IndexVersion: currentIndexVersion,
		}
		query := model.ValidatorSummary{
			Summary: summary,

			Address: rawSeqSummaryItem.Address,
		}

		existingValidatorSummary, err := uc.db.GetValidators().ValidatorSummary.Find(&query)
		if err != nil {
			if err == psql.ErrNotFound {
				validatorSummary := model.ValidatorSummary{
					Summary: summary,

					Address: rawSeqSummaryItem.Address,

					ScoreAvg: rawSeqSummaryItem.ScoreAvg,
					ScoreMin: rawSeqSummaryItem.ScoreMin,
					ScoreMax: rawSeqSummaryItem.ScoreMax,
					SignedAvg: rawSeqSummaryItem.SignedAvg,
					SignedMin: rawSeqSummaryItem.SignedMin,
					SignedMax: rawSeqSummaryItem.SignedMax,
				}

				if err := uc.db.GetValidators().ValidatorSummary.Create(&validatorSummary); err != nil {
					return err
				}
				newModels = append(newModels, validatorSummary)
			} else {
				return err
			}
		} else {
			existingValidatorSummary.ScoreAvg = rawSeqSummaryItem.ScoreAvg
			existingValidatorSummary.ScoreMin = rawSeqSummaryItem.ScoreMin
			existingValidatorSummary.ScoreMax = rawSeqSummaryItem.ScoreMax
			existingValidatorSummary.SignedAvg = rawSeqSummaryItem.SignedAvg
			existingValidatorSummary.SignedMin = rawSeqSummaryItem.SignedMin
			existingValidatorSummary.SignedMax = rawSeqSummaryItem.SignedMax

			if err := uc.db.GetValidators().ValidatorSummary.Save(existingValidatorSummary); err != nil {
				return err
			}
			existingModels = append(existingModels, *existingValidatorSummary)
		}
	}

	logger.Info(fmt.Sprintf("validator sequences summarized [created=%d] [updated=%d]", len(newModels), len(existingModels)))

	return nil
}

func (uc *summarizeUseCase) summarizeValidatorGroupSeq(interval types.SummaryInterval, currentIndexVersion int64) error {
	logger.Info(fmt.Sprintf("summarizing validator group sequences... [interval=%s]", interval))

	activityPeriods, err := uc.db.GetValidatorGroups().ValidatorGroupSummary.FindActivityPeriods(interval, currentIndexVersion)
	if err != nil {
		return err
	}

	rawSeqSummaryItems, err := uc.db.GetValidatorGroups().ValidatorGroupSeq.Summarize(interval, activityPeriods)
	if err != nil {
		return err
	}

	var newModels []model.ValidatorGroupSummary
	var existingModels []model.ValidatorGroupSummary
	for _, rawSeqSummaryItem := range rawSeqSummaryItems {
		summary := &model.Summary{
			TimeInterval: interval,
			TimeBucket:   rawSeqSummaryItem.TimeBucket,
			IndexVersion: currentIndexVersion,
		}
		query := model.ValidatorGroupSummary{
			Summary: summary,

			Address: rawSeqSummaryItem.Address,
		}

		existingValidatorGroupSummary, err := uc.db.GetValidatorGroups().ValidatorGroupSummary.Find(&query)
		if err != nil {
			if err == psql.ErrNotFound {
				validatorGroupSummary := model.ValidatorGroupSummary{
					Summary: summary,

					Address: rawSeqSummaryItem.Address,

					CommissionAvg: rawSeqSummaryItem.CommissionAvg,
					CommissionMin: rawSeqSummaryItem.CommissionMin,
					CommissionMax: rawSeqSummaryItem.CommissionMax,
					ActiveVotesAvg: rawSeqSummaryItem.ActiveVotesAvg,
					ActiveVotesMin: rawSeqSummaryItem.ActiveVotesMin,
					ActiveVotesMax: rawSeqSummaryItem.ActiveVotesMax,
					ActiveVoteUnitsAvg: rawSeqSummaryItem.ActiveVoteUnitsAvg,
					ActiveVoteUnitsMin: rawSeqSummaryItem.ActiveVoteUnitsMin,
					ActiveVoteUnitsMax: rawSeqSummaryItem.ActiveVoteUnitsMax,
					PendingVotesAvg: rawSeqSummaryItem.PendingVotesAvg,
					PendingVotesMin: rawSeqSummaryItem.PendingVotesMin,
					PendingVotesMax: rawSeqSummaryItem.PendingVotesMax,
				}

				if err := uc.db.GetValidatorGroups().ValidatorGroupSummary.Create(&validatorGroupSummary); err != nil {
					return err
				}
				newModels = append(newModels, validatorGroupSummary)
			} else {
				return err
			}
		} else {
			existingValidatorGroupSummary.CommissionAvg = rawSeqSummaryItem.CommissionAvg
			existingValidatorGroupSummary.CommissionMin = rawSeqSummaryItem.CommissionMin
			existingValidatorGroupSummary.CommissionMax = rawSeqSummaryItem.CommissionMax
			existingValidatorGroupSummary.ActiveVotesAvg = rawSeqSummaryItem.ActiveVotesAvg
			existingValidatorGroupSummary.ActiveVotesMin = rawSeqSummaryItem.ActiveVotesMin
			existingValidatorGroupSummary.ActiveVotesMax = rawSeqSummaryItem.ActiveVotesMax
			existingValidatorGroupSummary.ActiveVoteUnitsAvg = rawSeqSummaryItem.ActiveVoteUnitsAvg
			existingValidatorGroupSummary.ActiveVoteUnitsMin = rawSeqSummaryItem.ActiveVoteUnitsMin
			existingValidatorGroupSummary.ActiveVoteUnitsMax = rawSeqSummaryItem.ActiveVoteUnitsMax
			existingValidatorGroupSummary.PendingVotesAvg = rawSeqSummaryItem.PendingVotesAvg
			existingValidatorGroupSummary.PendingVotesMin = rawSeqSummaryItem.PendingVotesMin
			existingValidatorGroupSummary.PendingVotesMax = rawSeqSummaryItem.PendingVotesMax

			if err := uc.db.GetValidatorGroups().ValidatorGroupSummary.Save(existingValidatorGroupSummary); err != nil {
				return err
			}
			existingModels = append(existingModels, *existingValidatorGroupSummary)
		}
	}

	logger.Info(fmt.Sprintf("validator group sequences summarized [created=%d] [updated=%d]", len(newModels), len(existingModels)))

	return nil
}