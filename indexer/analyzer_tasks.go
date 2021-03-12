package indexer

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/figment-networks/celo-indexer/metric"
	"github.com/figment-networks/celo-indexer/store"
	"github.com/figment-networks/celo-indexer/store/psql"

	"github.com/figment-networks/celo-indexer/config"
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/types"
	"github.com/figment-networks/celo-indexer/utils/logger"
	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/pkg/errors"
)

const (
	TaskNameSystemEventCreator = "SystemEventCreator"
)

var (
	ErrGroupRewardOutsideOfRange = errors.New("group reward is outside of specified buckets")

	MaxValidatorSequences int64 = 1000
	MissedForMaxThreshold int64 = 50
	MissedInRowThreshold  int64 = 50
)

// NewSystemEventCreatorTask creates system events
func NewSystemEventCreatorTask(cfg *config.Config, validatorSeqDb store.ValidatorSeq, accountActivitySeqDb store.AccountActivitySeq) *systemEventCreatorTask {
	return &systemEventCreatorTask{
		validatorSeqDb:       validatorSeqDb,
		accountActivitySeqDb: accountActivitySeqDb,
		cfg:                  cfg,
	}
}

type systemEventCreatorTask struct {
	validatorSeqDb       store.ValidatorSeq
	accountActivitySeqDb store.AccountActivitySeq

	cfg *config.Config
}

type systemEventRawData map[string]interface{}

func (t *systemEventCreatorTask) GetName() string {
	return TaskNameSystemEventCreator
}

func (t *systemEventCreatorTask) Run(ctx context.Context, p pipeline.Payload) error {
	defer metric.LogIndexerTaskDuration(time.Now(), t.GetName())

	payload := p.(*payload)

	logger.Info(fmt.Sprintf("running indexer task [stage=%s] [task=%s] [height=%d]", "Analyzer", t.GetName(), payload.CurrentHeight))

	currEpochAccountActivitySequences := payload.AccountActivitySequences
	prevEpochAccountActivitySequences, err := t.getPrevEpochAccountActivitySequences(payload)
	if err != nil {
		return err
	}

	currHeightValidatorSequences := payload.ValidatorSequences
	prevHeightValidatorSequences, err := t.getPrevHeightValidatorSequences(payload)
	if err != nil {
		return err
	}

	valueChangeSystemEvents, err := t.getValueChangeForAccountActivity(payload.HeightMeta, currEpochAccountActivitySequences, prevEpochAccountActivitySequences)
	if err != nil {
		return err
	}
	payload.SystemEvents = append(payload.SystemEvents, valueChangeSystemEvents...)

	activeSetPresenceChangeSystemEvents, err := t.getActiveSetPresenceChangeSystemEvents(currHeightValidatorSequences, prevHeightValidatorSequences)
	if err != nil {
		return err
	}
	payload.SystemEvents = append(payload.SystemEvents, activeSetPresenceChangeSystemEvents...)

	missedBlocksSystemEvents, err := t.getMissedBlocksSystemEvents(currHeightValidatorSequences)
	if err != nil {
		return err
	}
	payload.SystemEvents = append(payload.SystemEvents, missedBlocksSystemEvents...)

	return nil
}

func (t *systemEventCreatorTask) getPrevEpochAccountActivitySequences(payload *payload) ([]model.AccountActivitySeq, error) {
	epochSize := payload.HeightMeta.EpochSize
	prevEpochHeight := payload.CurrentHeight - *epochSize

	var prevEpochAccountActivitySequences []model.AccountActivitySeq
	if payload.CurrentHeight > t.cfg.FirstBlockHeight {
		var err error
		prevEpochAccountActivitySequences, err = t.accountActivitySeqDb.FindByHeight(prevEpochHeight)
		if err != nil {
			if err != psql.ErrNotFound {
				return nil, err
			}
		}
	}
	return prevEpochAccountActivitySequences, nil
}

func (t *systemEventCreatorTask) getPrevHeightValidatorSequences(payload *payload) ([]model.ValidatorSeq, error) {
	var prevHeightValidatorSequences []model.ValidatorSeq
	if payload.CurrentHeight > t.cfg.FirstBlockHeight {
		var err error
		prevHeightValidatorSequences, err = t.validatorSeqDb.FindByHeight(payload.CurrentHeight - 1)
		if err != nil {
			if err != psql.ErrNotFound {
				return nil, err
			}
		}
	}
	return prevHeightValidatorSequences, nil
}

func (t *systemEventCreatorTask) getMissedBlocksSystemEvents(currHeightValidatorSequences []model.ValidatorSeq) ([]model.SystemEvent, error) {
	var systemEvents []model.SystemEvent
	for _, validatorSequence := range currHeightValidatorSequences {
		// When current height validator has validated the block no need to check last records
		if t.isValidated(validatorSequence) {
			return systemEvents, nil
		}

		lastValidatorSequencesForAddress, err := t.validatorSeqDb.FindLastByAddress(validatorSequence.Address, MaxValidatorSequences)
		if err != nil {
			if err == psql.ErrNotFound {
				return systemEvents, nil
			} else {
				return nil, err
			}
		} else {
			var validatorSequencesToCheck []model.ValidatorSeq
			validatorSequencesToCheck = append([]model.ValidatorSeq{validatorSequence}, lastValidatorSequencesForAddress...)
			totalMissedCount := t.getTotalMissed(validatorSequencesToCheck)

			logger.Debug(fmt.Sprintf("total missed blocks for last %d blocks for address %s: %d", MaxValidatorSequences, validatorSequence.Address, totalMissedCount))

			if totalMissedCount == MissedForMaxThreshold {
				newSystemEvent, err := t.newSystemEvent(validatorSequence.Sequence, validatorSequence.Address, model.SystemEventMissedNofM, systemEventRawData{
					"threshold":               MissedForMaxThreshold,
					"max_validator_sequences": MaxValidatorSequences,
				})
				if err != nil {
					return nil, err
				}

				systemEvents = append(systemEvents, *newSystemEvent)
			}

			missedInRowCount := t.getMissedInRow(validatorSequencesToCheck, MissedInRowThreshold)

			logger.Debug(fmt.Sprintf("total missed blocks in a row for address %s: %d", validatorSequence.Address, missedInRowCount))

			if missedInRowCount == MissedInRowThreshold {
				newSystemEvent, err := t.newSystemEvent(validatorSequence.Sequence, validatorSequence.Address, model.SystemEventMissedNConsecutive, systemEventRawData{
					"threshold": MissedInRowThreshold,
				})
				if err != nil {
					return nil, err
				}

				systemEvents = append(systemEvents, *newSystemEvent)
			}
		}
	}
	return systemEvents, nil
}

// getTotalMissed get total missed count for given slice of validator sequences
func (t systemEventCreatorTask) getTotalMissed(validatorSequences []model.ValidatorSeq) int64 {
	var totalMissedCount int64 = 0
	for _, validatorSequence := range validatorSequences {
		if t.isNotValidated(validatorSequence) {
			totalMissedCount += 1
		}
	}

	return totalMissedCount
}

// getMissedInRow get number of validator sequences missed in the row
func (t systemEventCreatorTask) getMissedInRow(validatorSequences []model.ValidatorSeq, limit int64) int64 {
	if int64(len(validatorSequences)) > MissedInRowThreshold {
		validatorSequences = validatorSequences[:limit]
	}

	var missedInRowCount int64 = 0
	prevValidated := false
	for _, validatorSequence := range validatorSequences {
		if t.isNotValidated(validatorSequence) {
			if !prevValidated {
				missedInRowCount += 1
			}
			prevValidated = false
		} else {
			prevValidated = true
		}
	}

	return missedInRowCount
}

// isNotValidated check if validator has validated the block at height
func (t systemEventCreatorTask) isNotValidated(validatorSequence model.ValidatorSeq) bool {
	return validatorSequence.Signed != nil && !*validatorSequence.Signed
}

func (t systemEventCreatorTask) isValidated(validatorSequence model.ValidatorSeq) bool {
	return !t.isNotValidated(validatorSequence)
}

func (t *systemEventCreatorTask) getActiveSetPresenceChangeSystemEvents(currHeightValidatorSequences []model.ValidatorSeq, prevHeightValidatorSequences []model.ValidatorSeq) ([]model.SystemEvent, error) {
	var systemEvents []model.SystemEvent

	for _, currentValidatorSequence := range currHeightValidatorSequences {
		joined := true
		for _, prevValidatorSequence := range prevHeightValidatorSequences {
			if currentValidatorSequence.Address == prevValidatorSequence.Address {
				joined = false
				break
			}
		}

		if joined {
			logger.Debug(fmt.Sprintf("address %s joined active set", currentValidatorSequence.Address))

			newSystemEvent, err := t.newSystemEvent(currentValidatorSequence.Sequence, currentValidatorSequence.Address, model.SystemEventJoinedActiveSet, systemEventRawData{})
			if err != nil {
				return nil, err
			}

			systemEvents = append(systemEvents, *newSystemEvent)
		}
	}

	for _, prevValidatorSequence := range prevHeightValidatorSequences {
		left := true
		for _, currentValidatorSequence := range currHeightValidatorSequences {
			if prevValidatorSequence.Address == currentValidatorSequence.Address {
				left = false
				break
			}
		}

		if left {
			logger.Debug(fmt.Sprintf("address %s joined active set", prevValidatorSequence.Address))

			newSystemEvent, err := t.newSystemEvent(prevValidatorSequence.Sequence, prevValidatorSequence.Address, model.SystemEventLeftActiveSet, systemEventRawData{})
			if err != nil {
				return nil, err
			}

			systemEvents = append(systemEvents, *newSystemEvent)
		}
	}

	return systemEvents, nil
}

func (t *systemEventCreatorTask) getValueChangeForAccountActivity(heightMeta HeightMeta, currHeightItems []model.AccountActivitySeq, prevHeightItems []model.AccountActivitySeq) ([]model.SystemEvent, error) {
	return t.getValueChangeForAccountActivityByKind(heightMeta, currHeightItems, prevHeightItems, OperationTypeValidatorEpochPaymentDistributedForGroup)
}

func (t *systemEventCreatorTask) getValueChangeForAccountActivityByKind(heightMeta HeightMeta, currHeightItems []model.AccountActivitySeq, prevHeightItems []model.AccountActivitySeq, kind string) ([]model.SystemEvent, error) {
	filteredCurrHeightItems := t.filterAccountActivitiesByKind(currHeightItems, kind)
	filteredPrevHeightItems := t.filterAccountActivitiesByKind(prevHeightItems, kind)

	aggregatedCurrHeightItems, err := t.aggregateAccountActivitiesByAddress(filteredCurrHeightItems)
	if err != nil {
		return nil, err
	}
	aggregatedPrevHeightItems, err := t.aggregateAccountActivitiesByAddress(filteredPrevHeightItems)
	if err != nil {
		return nil, err
	}

	var systemEvents []model.SystemEvent

	for currAddress, currAmount := range aggregatedCurrHeightItems {
		var previousAmount int64
		found := false
		for prevAddress, prevAmount := range aggregatedPrevHeightItems {
			if currAddress == prevAddress {
				previousAmount = prevAmount.Int64()
				found = true
				break
			}
		}

		if !found {
			// Is in current but not in previous
			previousAmount = 0
		}

		// Create system events for both existing and missing in previous account activities
		newSystemEvent, err := t.getGroupRewardDistributedChange(currAmount.Int64(), previousAmount)
		if err != nil {
			if err != ErrGroupRewardOutsideOfRange {
				return nil, err
			}
		} else {
			newSystemEvent.Height = heightMeta.Height
			newSystemEvent.Time = *heightMeta.Time
			newSystemEvent.Actor = currAddress
			logger.Debug(fmt.Sprintf("group reward change for address %s occured [kind=%s]", currAddress, newSystemEvent.Kind))
			systemEvents = append(systemEvents, *newSystemEvent)
		}
	}

	for prevAddress, prevAmount := range aggregatedPrevHeightItems {
		found := false
		for currAddress := range aggregatedCurrHeightItems {
			if prevAddress == currAddress {
				found = true
				break
			}
		}

		if !found {
			// Is in previous but not in current
			// We only have to create system events for missing records
			currentAmount := int64(0)
			newSystemEvent, err := t.getGroupRewardDistributedChange(currentAmount, prevAmount.Int64())
			if err != nil {
				if err != ErrGroupRewardOutsideOfRange {
					return nil, err
				}
			} else {
				newSystemEvent.Height = heightMeta.Height
				newSystemEvent.Time = *heightMeta.Time
				newSystemEvent.Actor = prevAddress
				logger.Debug(fmt.Sprintf("group reward change for address %s occured [kind=%s]", prevAddress, newSystemEvent.Kind))
				systemEvents = append(systemEvents, *newSystemEvent)
			}
		}
	}

	return systemEvents, nil
}

func (t *systemEventCreatorTask) getGroupRewardDistributedChange(currValue int64, prevValue int64) (*model.SystemEvent, error) {
	roundedChangeRate := t.getRoundedChangeRate(currValue, prevValue)
	roundedAbsChangeRate := math.Abs(roundedChangeRate)

	var kind model.SystemEventKind
	if roundedAbsChangeRate >= 0.1 && roundedAbsChangeRate < 1 {
		kind = model.SystemEventGroupRewardChange1
	} else if roundedAbsChangeRate >= 1 && roundedAbsChangeRate < 10 {
		kind = model.SystemEventGroupRewardChange2
	} else if roundedAbsChangeRate >= 10 {
		kind = model.SystemEventGroupRewardChange3
	} else {
		return nil, ErrGroupRewardOutsideOfRange
	}

	data := systemEventRawData{
		"before": prevValue,
		"after":  currValue,
		"change": roundedChangeRate,
	}
	marshaledData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return &model.SystemEvent{
		Kind: kind,
		Data: types.Jsonb{RawMessage: marshaledData},
	}, nil
}

func (t *systemEventCreatorTask) filterAccountActivitiesByKind(accountActivitySeqs []model.AccountActivitySeq, kind string) []model.AccountActivitySeq {
	var filtered []model.AccountActivitySeq
	for _, seq := range accountActivitySeqs {
		if seq.Kind == kind {
			filtered = append(filtered, seq)
		}
	}
	return filtered
}

func (t *systemEventCreatorTask) aggregateAccountActivitiesByAddress(accountActivitySeqs []model.AccountActivitySeq) (map[string]types.Quantity, error) {
	aggregated := make(map[string]types.Quantity)

	for _, seq := range accountActivitySeqs {
		currAmount, ok := aggregated[seq.Address]
		if ok {
			if err := currAmount.Add(seq.Amount); err != nil {
				return nil, err
			}
		} else {
			aggregated[seq.Address] = seq.Amount
		}
	}
	return aggregated, nil
}

func (t *systemEventCreatorTask) getRoundedChangeRate(currValue int64, prevValue int64) float64 {
	var changeRate float64

	if prevValue == 0 {
		changeRate = float64(currValue)
	} else {
		changeRate = (float64(1) - (float64(currValue) / float64(prevValue))) * 100
	}

	roundedChangeRate := math.Round(changeRate/0.1) * 0.1
	return roundedChangeRate
}

func (t *systemEventCreatorTask) newSystemEvent(seq *model.Sequence, actor string, kind model.SystemEventKind, data map[string]interface{}) (*model.SystemEvent, error) {
	marshaledData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return &model.SystemEvent{
		Height: seq.Height,
		Time:   seq.Time,
		Actor:  actor,
		Kind:   kind,
		Data:   types.Jsonb{RawMessage: marshaledData},
	}, nil
}
