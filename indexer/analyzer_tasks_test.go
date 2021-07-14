package indexer

import (
	"testing"
	"time"

	"github.com/figment-networks/celo-indexer/config"
	mock "github.com/figment-networks/celo-indexer/mock/store"
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/store/psql"
	"github.com/figment-networks/celo-indexer/types"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
)

const (
	testAccountAddress        = "test_acct_address"
	testValidatorAddress      = "test_address"
	testValidatorGroupAddress = "test_validator_group_address"
	testHeight                = 17
)

var (
	ErrCouldNotFindByAddress = errors.New("could not find test")

	testCfg = &config.Config{
		FirstBlockHeight: 1,
	}
)

func TestSystemEventCreatorTask_GroupRewardDistributedChangeSystemEvents(t *testing.T) {
	currSyncable := &model.Syncable{
		Height: 20,
		Time:   types.NewTimeFromTime(time.Date(2020, 11, 10, 23, 0, 0, 0, time.UTC)),
	}
	prevSyncable := &model.Syncable{
		Height: 19,
		Time:   types.NewTimeFromTime(time.Date(2020, 11, 10, 22, 0, 0, 0, time.UTC)),
	}

	currSeq := &model.Sequence{
		Height: currSyncable.Height,
		Time:   *currSyncable.Time,
	}
	prevSeq := &model.Sequence{
		Height: prevSyncable.Height,
		Time:   *prevSyncable.Time,
	}

	tests := []struct {
		description           string
		groupRewardChangeRate float64
		expectedCount         int
		expectedKind          model.SystemEventKind
	}{
		{"returns no system events when group reward haven't changed", 0, 0, ""},
		{"returns no system events when group reward change smaller than 0.1", 0.09, 0, ""},
		{"returns one groupRewardChange1 system event when group reward change is 0.1", 0.1, 2, model.SystemEventGroupRewardChange1},
		{"returns one groupRewardChange1 system events when group reward change is 0.9", 0.9, 2, model.SystemEventGroupRewardChange1},
		{"returns one groupRewardChange2 system events when group reward change is 1", 1, 2, model.SystemEventGroupRewardChange2},
		{"returns one groupRewardChange2 system events when group reward change is 9", 9, 2, model.SystemEventGroupRewardChange2},
		{"returns one groupRewardChange3 system events when group reward change is 10", 10, 2, model.SystemEventGroupRewardChange3},
		{"returns one groupRewardChange3 system events when group reward change is 100", 100, 2, model.SystemEventGroupRewardChange3},
		{"returns one groupRewardChange3 system events when group reward change is 200", 200, 2, model.SystemEventGroupRewardChange3},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			var groupRewardBefore int64 = 1000
			groupRewardAfter := float64(groupRewardBefore) + (float64(groupRewardBefore) * tt.groupRewardChangeRate / 100)

			prevHeightAccountActivitySequences := []model.AccountActivitySeq{
				{
					Sequence: prevSeq,
					Address:  testAccountAddress,
					Kind:     OperationTypeValidatorEpochPaymentDistributedForGroup,
					Amount:   types.NewQuantityFromInt64(groupRewardBefore),
				},
				{
					Sequence: currSeq,
					Address:  testValidatorAddress,
					Kind:     OperationTypeValidatorEpochPaymentDistributedForGroup,
					Amount:   types.NewQuantityFromInt64(groupRewardBefore),
				},
			}
			currHeightAccountActivitySequences := []model.AccountActivitySeq{
				{
					Sequence: prevSeq,
					Address:  testAccountAddress,
					Kind:     OperationTypeValidatorEpochPaymentDistributedForGroup,
					Amount:   types.NewQuantityFromInt64(int64(groupRewardAfter)),
				},
				{
					Sequence: currSeq,
					Address:  testValidatorAddress,
					Kind:     OperationTypeValidatorEpochPaymentDistributedForGroup,
					Amount:   types.NewQuantityFromInt64(int64(groupRewardAfter)),
				},
			}

			heightMeta := HeightMeta{
				Height: currSyncable.Height,
				Time:   currSyncable.Time,
			}
			task := NewSystemEventCreatorTask(testCfg, nil, nil, nil)
			createdSystemEvents, _ := task.getValueChangeForAccountActivityByKind(heightMeta, currHeightAccountActivitySequences, prevHeightAccountActivitySequences, OperationTypeValidatorEpochPaymentDistributedForGroup)

			if len(createdSystemEvents) != tt.expectedCount {
				t.Errorf("unexpected system event count, want %v; got %v", tt.expectedCount, len(createdSystemEvents))
				return
			}

			if len(createdSystemEvents) > 0 && createdSystemEvents[0].Kind != tt.expectedKind {
				t.Errorf("unexpected system event kind, want %v; got %v", tt.expectedKind, createdSystemEvents[0].Kind)
			}
		})
	}
}

func TestSystemEventCreatorTask_getActiveSetPresenceChangeSystemEvents(t *testing.T) {
	currSyncable := &model.Syncable{
		Height: 20,
		Time:   types.NewTimeFromTime(time.Date(2020, 11, 10, 23, 0, 0, 0, time.UTC)),
	}
	prevSyncable := &model.Syncable{
		Height: 19,
		Time:   types.NewTimeFromTime(time.Date(2020, 11, 10, 22, 0, 0, 0, time.UTC)),
	}

	currSeq := &model.Sequence{
		Height: currSyncable.Height,
		Time:   *currSyncable.Time,
	}
	prevSeq := &model.Sequence{
		Height: prevSyncable.Height,
		Time:   *prevSyncable.Time,
	}

	tests := []struct {
		description   string
		prevSeqs      []model.ValidatorSeq
		currSeqs      []model.ValidatorSeq
		expectedCount int
		expectedKinds []model.SystemEventKind
	}{
		{
			description: "returns no system events when validator is both in prev and current lists",
			prevSeqs: []model.ValidatorSeq{
				{
					Sequence: prevSeq,
					Address:  testValidatorAddress,
				},
			},
			currSeqs: []model.ValidatorSeq{
				{
					Sequence: currSeq,
					Address:  testValidatorAddress,
				},
			},
			expectedCount: 0,
		},
		{
			description:   "returns no system events when validator is not in any list",
			prevSeqs:      []model.ValidatorSeq{},
			currSeqs:      []model.ValidatorSeq{},
			expectedCount: 0,
		},
		{
			description: "returns one joined_set system events when validator is not in prev lists and is in current list",
			prevSeqs:    []model.ValidatorSeq{},
			currSeqs: []model.ValidatorSeq{{
				Sequence: currSeq,
				Address:  testValidatorAddress,
			}},
			expectedCount: 1,
			expectedKinds: []model.SystemEventKind{model.SystemEventJoinedActiveSet},
		},
		{
			description: "returns one left_set system events when validator is in prev set and not in current set",
			prevSeqs: []model.ValidatorSeq{
				{
					Sequence: prevSeq,
					Address:  testValidatorAddress,
				},
			},
			currSeqs:      []model.ValidatorSeq{},
			expectedCount: 1,
			expectedKinds: []model.SystemEventKind{model.SystemEventLeftActiveSet},
		},
		{
			description: "returns 2 joined_set system events when validators are not in prev but are in current lists",
			prevSeqs:    []model.ValidatorSeq{},
			currSeqs: []model.ValidatorSeq{
				{
					Sequence: currSeq,
					Address:  testValidatorAddress,
				},
				{
					Sequence: currSeq,
					Address:  "testValidatorAddress2",
				},
			},
			expectedCount: 2,
			expectedKinds: []model.SystemEventKind{model.SystemEventJoinedActiveSet, model.SystemEventJoinedActiveSet},
		},
		{
			description: "returns 2 left_set system events when validators are in prev but are not in current lists",
			prevSeqs: []model.ValidatorSeq{
				{
					Sequence: prevSeq,
					Address:  testValidatorAddress,
				},
				{
					Sequence: prevSeq,
					Address:  "testValidatorAddress2",
				},
			},
			currSeqs:      []model.ValidatorSeq{},
			expectedCount: 2,
			expectedKinds: []model.SystemEventKind{model.SystemEventLeftActiveSet, model.SystemEventLeftActiveSet},
		},
		{
			description: "returns left and joined set events",
			prevSeqs: []model.ValidatorSeq{
				{
					Sequence: prevSeq,
					Address:  testValidatorAddress,
				},
				{
					Sequence: prevSeq,
					Address:  "testValidatorAddress2",
				},
			},
			currSeqs: []model.ValidatorSeq{
				{
					Sequence: currSeq,
					Address:  testValidatorAddress,
				},
				{
					Sequence: currSeq,
					Address:  "testValidatorAddress3",
				},
			},
			expectedCount: 2,
			expectedKinds: []model.SystemEventKind{model.SystemEventJoinedActiveSet, model.SystemEventLeftActiveSet},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			task := NewSystemEventCreatorTask(testCfg, nil, nil, nil)
			createdSystemEvents, _ := task.getActiveSetPresenceChangeSystemEvents(tt.currSeqs, tt.prevSeqs)

			if len(createdSystemEvents) != tt.expectedCount {
				t.Errorf("unexpected system event count, want %v; got %v", tt.expectedCount, len(createdSystemEvents))
				return
			}

			for i, kind := range tt.expectedKinds {
				if len(createdSystemEvents) > 0 && createdSystemEvents[i].Kind != kind {
					t.Errorf("unexpected system event kind, want %v; got %v", kind, createdSystemEvents[i].Kind)
				}
			}
		})
	}
}

func TestSystemEventCreatorTask_getMissedBlocksSystemEventsForValidatorSequences(t *testing.T) {
	tests := []struct {
		description           string
		maxValidatorSequences int64
		missedInRowThreshold  int64
		missedForMaxThreshold int64
		prevHeightList        []model.ValidatorSeq
		currHeightList        []model.ValidatorSeq
		lastForValidatorList  [][]model.ValidatorSeq
		errs                  []error
		expectedCount         int
		expectedKinds         []model.SystemEventKind
		expectedErr           error
	}{
		{
			description:           "returns no system events when validator does not have any previous sequences in db",
			maxValidatorSequences: 5,
			missedInRowThreshold:  2,
			missedForMaxThreshold: 2,
			prevHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, true),
			},
			currHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, false),
			},
			lastForValidatorList: [][]model.ValidatorSeq{
				{},
			},
			expectedCount: 0,
		},
		{
			description:           "returns no system events when validator does not have any missed blocks in a row",
			maxValidatorSequences: 5,
			missedInRowThreshold:  2,
			missedForMaxThreshold: 2,
			prevHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, true),
			},
			currHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, false),
			},
			lastForValidatorList: [][]model.ValidatorSeq{
				{
					newValidatorSeq(testValidatorAddress, 1000, true),
					newValidatorSeq(testValidatorAddress, 1000, true),
					newValidatorSeq(testValidatorAddress, 1000, true),
					newValidatorSeq(testValidatorAddress, 1000, true),
					newValidatorSeq(testValidatorAddress, 1000, true),
				},
			},
			expectedCount: 0,
		},
		{
			description:           "returns no system events when validator missed 2 blocks in a row",
			maxValidatorSequences: 5,
			missedInRowThreshold:  3,
			missedForMaxThreshold: 5,
			prevHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, true),
			},
			currHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, false),
			},
			lastForValidatorList: [][]model.ValidatorSeq{
				{
					newValidatorSeq(testValidatorAddress, 1000, false),
					newValidatorSeq(testValidatorAddress, 1000, true),
					newValidatorSeq(testValidatorAddress, 1000, true),
					newValidatorSeq(testValidatorAddress, 1000, true),
				},
			},
			expectedCount: 0,
		},
		{
			description:           "returns one missed_n_consecutive system events when validator missed >= 3 blocks in a row",
			maxValidatorSequences: 5,
			missedInRowThreshold:  3,
			missedForMaxThreshold: 5,
			prevHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, true),
			},
			currHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, false),
			},
			lastForValidatorList: [][]model.ValidatorSeq{
				{
					newValidatorSeq(testValidatorAddress, 1000, false),
					newValidatorSeq(testValidatorAddress, 1000, false),
					newValidatorSeq(testValidatorAddress, 1000, true),
					newValidatorSeq(testValidatorAddress, 1000, true),
				},
			},
			expectedCount: 1,
			expectedKinds: []model.SystemEventKind{model.SystemEventMissedNConsecutive},
		},
		{
			description:           "returns no missed_n_consecutive system events when validator missed >= 3 blocks in a row in the past but current is validated",
			maxValidatorSequences: 5,
			missedInRowThreshold:  3,
			missedForMaxThreshold: 5,
			prevHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, true),
			},
			currHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, true),
			},
			lastForValidatorList: [][]model.ValidatorSeq{
				{
					newValidatorSeq(testValidatorAddress, 1000, false),
					newValidatorSeq(testValidatorAddress, 1000, false),
					newValidatorSeq(testValidatorAddress, 1000, false),
					newValidatorSeq(testValidatorAddress, 1000, true),
				},
			},
			expectedCount: 0,
		},
		{
			description:           "returns one missed_n_of_m system events when validator missed 3 blocks",
			maxValidatorSequences: 5,
			missedInRowThreshold:  50,
			missedForMaxThreshold: 3,
			prevHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, true),
			},
			currHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, false),
			},
			lastForValidatorList: [][]model.ValidatorSeq{
				{
					newValidatorSeq(testValidatorAddress, 1000, false),
					newValidatorSeq(testValidatorAddress, 1000, false),
					newValidatorSeq(testValidatorAddress, 1000, true),
					newValidatorSeq(testValidatorAddress, 1000, true),
				},
			},
			expectedCount: 1,
			expectedKinds: []model.SystemEventKind{model.SystemEventMissedNofM},
		},
		{
			description:           "returns one missed_n_of_m system events when validator missed 3 blocks and max < last list",
			maxValidatorSequences: 3,
			missedInRowThreshold:  50,
			missedForMaxThreshold: 3,
			prevHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, true),
			},
			currHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, false),
			},
			lastForValidatorList: [][]model.ValidatorSeq{
				{
					newValidatorSeq(testValidatorAddress, 1000, false),
					newValidatorSeq(testValidatorAddress, 1000, false),
					newValidatorSeq(testValidatorAddress, 1000, true),
					newValidatorSeq(testValidatorAddress, 1000, true),
					newValidatorSeq(testValidatorAddress, 1000, true),
				},
			},
			expectedCount: 1,
			expectedKinds: []model.SystemEventKind{model.SystemEventMissedNofM},
		},
		{
			description:           "returns no missed_n_of_m system events when count of recent not validated > maxValidatorSequences",
			maxValidatorSequences: 5,
			missedInRowThreshold:  50,
			missedForMaxThreshold: 3,
			prevHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, true),
			},
			currHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, false),
			},
			lastForValidatorList: [][]model.ValidatorSeq{
				{
					newValidatorSeq(testValidatorAddress, 1000, false),
					newValidatorSeq(testValidatorAddress, 1000, false),
					newValidatorSeq(testValidatorAddress, 1000, false),
					newValidatorSeq(testValidatorAddress, 1000, false),
					newValidatorSeq(testValidatorAddress, 1000, false),
				},
			},
			expectedCount: 0,
		},
		{
			description:           "returns no missed_n_of_m system events when current is validated",
			maxValidatorSequences: 5,
			missedInRowThreshold:  50,
			missedForMaxThreshold: 3,
			prevHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, true),
			},
			currHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, true),
			},
			lastForValidatorList: [][]model.ValidatorSeq{
				{
					newValidatorSeq(testValidatorAddress, 1000, false),
					newValidatorSeq(testValidatorAddress, 1000, false),
					newValidatorSeq(testValidatorAddress, 1000, false),
					newValidatorSeq(testValidatorAddress, 1000, false),
					newValidatorSeq(testValidatorAddress, 1000, false),
				},
			},
			expectedCount: 0,
		},
		{
			description:           "returns error when first call to FindLastByAddress fails",
			maxValidatorSequences: 3,
			missedInRowThreshold:  50,
			missedForMaxThreshold: 3,
			prevHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, true),
			},
			currHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, false),
			},
			lastForValidatorList: [][]model.ValidatorSeq{
				nil,
			},
			errs:          []error{ErrCouldNotFindByAddress},
			expectedCount: 0,
			expectedErr:   ErrCouldNotFindByAddress,
		},
		{
			description:           "returns error when second call to FindLastByAddress fails",
			maxValidatorSequences: 5,
			missedInRowThreshold:  3,
			missedForMaxThreshold: 5,
			prevHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, true),
				newValidatorSeq("address1", 1000, false),
			},
			currHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, false),
				newValidatorSeq("address1", 1000, false),
			},
			lastForValidatorList: [][]model.ValidatorSeq{
				{
					newValidatorSeq(testValidatorAddress, 1000, false),
					newValidatorSeq(testValidatorAddress, 1000, false),
					newValidatorSeq(testValidatorAddress, 1000, false),
					newValidatorSeq(testValidatorAddress, 1000, true),
				},
				nil,
			},
			errs:          []error{nil, ErrCouldNotFindByAddress},
			expectedCount: 0,
			expectedErr:   ErrCouldNotFindByAddress,
		},
		{
			description:           "returns partial system events when second call to FindLastByAddress fails with ErrNotFound",
			maxValidatorSequences: 3,
			missedInRowThreshold:  50,
			missedForMaxThreshold: 3,
			prevHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, true),
				newValidatorSeq("address1", 1000, true),
			},
			currHeightList: []model.ValidatorSeq{
				newValidatorSeq(testValidatorAddress, 1000, false),
				newValidatorSeq("address1", 1000, false),
			},
			lastForValidatorList: [][]model.ValidatorSeq{
				{
					newValidatorSeq(testValidatorAddress, 1000, false),
					newValidatorSeq(testValidatorAddress, 1000, false),
					newValidatorSeq(testValidatorAddress, 1000, true),
					newValidatorSeq(testValidatorAddress, 1000, true),
				},
				nil,
			},
			errs:          []error{nil, psql.ErrNotFound},
			expectedCount: 1,
			expectedKinds: []model.SystemEventKind{model.SystemEventMissedNofM},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			validatorSeqStoreMock := mock.NewMockValidatorSeq(ctrl)
			accountActivitySeqStoreMock := mock.NewMockAccountActivitySeq(ctrl)

			MaxValidatorSequences = tt.maxValidatorSequences
			MissedInRowThreshold = tt.missedInRowThreshold
			MissedForMaxThreshold = tt.missedForMaxThreshold

			var mockCalls []*gomock.Call
			for i, validatorSeqs := range tt.lastForValidatorList {
				validatorSeq := tt.currHeightList[i]
				if validatorSeq.Signed == nil || !*validatorSeq.Signed {
					call := validatorSeqStoreMock.EXPECT().FindLastByAddress(gomock.Any(), gomock.Any())

					if len(tt.errs) >= i+1 && tt.errs[i] != nil {
						call = call.Return(nil, tt.errs[i])
					} else {
						call = call.Return(validatorSeqs, nil)
					}

					mockCalls = append(mockCalls, call)
				}
			}
			gomock.InOrder(mockCalls...)

			task := NewSystemEventCreatorTask(testCfg, validatorSeqStoreMock, accountActivitySeqStoreMock, nil)
			createdSystemEvents, err := task.getMissedBlocksOfValidatorSequences(tt.currHeightList)
			if err == nil && tt.expectedErr != nil {
				t.Errorf("should return error")
				return
			}
			if err != nil && tt.expectedErr != err {
				t.Errorf("unexpected error, want %v; got %v", tt.expectedErr, err)
				return
			}

			if len(createdSystemEvents) != tt.expectedCount {
				t.Errorf("unexpected system event count, want %v; got %v", tt.expectedCount, len(createdSystemEvents))
				return
			}

			for i, kind := range tt.expectedKinds {
				if len(createdSystemEvents) > 0 && createdSystemEvents[i].Kind != kind {
					t.Errorf("unexpected system event kind, want %v; got %v", kind, createdSystemEvents[i].Kind)
				}
			}
		})
	}
}

func newValidatorSeq(address string, score int64, signed bool) model.ValidatorSeq {
	return model.ValidatorSeq{
		Sequence: &model.Sequence{
			Height: testHeight,
			Time:   *types.NewTimeFromTime(time.Now()),
		},
		Address: address,
		Score:   types.NewQuantityFromInt64(score),
		Signed:  &signed,
	}
}

func TestSystemEventCreatorTask_getMissedBlocksSystemEventsForValidatorGroupSequences(t *testing.T) {
	tests := []struct {
		description                string
		maxValidatorGroupSequences int64
		missedInRowThreshold       int64
		missedForMaxThreshold      int64
		currHeightList             []model.ValidatorGroupSeq
		lastForValidatorGroupList  [][]model.ValidatorGroupSeq
		errs                       []error
		expectedCount              int
		expectedKinds              []model.SystemEventKind
		expectedErr                error
	}{
		{
			description:                "returns no system events when validator group does not have any previous sequences in db",
			maxValidatorGroupSequences: 5,
			missedInRowThreshold:       2,
			missedForMaxThreshold:      2,
			currHeightList: []model.ValidatorGroupSeq{
				newValidatorGroupSeq(testValidatorGroupAddress, 10, 0),
			},
			lastForValidatorGroupList: [][]model.ValidatorGroupSeq{
				{},
			},
			expectedCount: 0,
		},
		{
			description:                "returns no system events when validator group does not have any missed blocks in a row",
			maxValidatorGroupSequences: 5,
			missedInRowThreshold:       2,
			missedForMaxThreshold:      2,
			currHeightList: []model.ValidatorGroupSeq{
				newValidatorGroupSeq(testValidatorGroupAddress, 1000, 0),
			},
			lastForValidatorGroupList: [][]model.ValidatorGroupSeq{
				{
					newValidatorGroupSeq(testValidatorGroupAddress, 1000, 0.1),
					newValidatorGroupSeq(testValidatorGroupAddress, 1000, 0.1),
					newValidatorGroupSeq(testValidatorGroupAddress, 1000, 0.1),
					newValidatorGroupSeq(testValidatorGroupAddress, 1000, 0.1),
					newValidatorGroupSeq(testValidatorGroupAddress, 1000, 0.1),
				},
			},
			expectedCount: 0,
		},
		{
			description:                "returns no system events when validator group missed 2 blocks in a row",
			maxValidatorGroupSequences: 5,
			missedInRowThreshold:       3,
			missedForMaxThreshold:      5,
			currHeightList: []model.ValidatorGroupSeq{
				newValidatorGroupSeq(testValidatorGroupAddress, 1000, 0),
			},
			lastForValidatorGroupList: [][]model.ValidatorGroupSeq{
				{
					newValidatorGroupSeq(testValidatorGroupAddress, 1000, 0),
					newValidatorGroupSeq(testValidatorGroupAddress, 1000, 0.1),
					newValidatorGroupSeq(testValidatorGroupAddress, 1000, 0.1),
					newValidatorGroupSeq(testValidatorGroupAddress, 1000, 0.1),
				},
			},
			expectedCount: 0,
		},
		{
			description:                "returns one missed_n_consecutive system events when validator group missed >= 3 blocks in a row",
			maxValidatorGroupSequences: 5,
			missedInRowThreshold:       3,
			missedForMaxThreshold:      5,
			currHeightList: []model.ValidatorGroupSeq{
				newValidatorGroupSeq(testValidatorGroupAddress, 1000, 0),
			},
			lastForValidatorGroupList: [][]model.ValidatorGroupSeq{
				{
					newValidatorGroupSeq(testValidatorGroupAddress, 1000, 0),
					newValidatorGroupSeq(testValidatorGroupAddress, 1000, 0),
					newValidatorGroupSeq(testValidatorGroupAddress, 1000, 0.1),
					newValidatorGroupSeq(testValidatorGroupAddress, 1000, 0.1),
				},
			},
			expectedCount: 1,
			expectedKinds: []model.SystemEventKind{model.SystemEventMissedNConsecutive},
		},
		{
			description:                "returns no missed_n_consecutive system events when validator group missed >= 3 blocks in a row in the past but current is validated",
			maxValidatorGroupSequences: 5,
			missedInRowThreshold:       3,
			missedForMaxThreshold:      5,
			currHeightList: []model.ValidatorGroupSeq{
				newValidatorGroupSeq(testValidatorAddress, 1000, 0.1),
			},
			lastForValidatorGroupList: [][]model.ValidatorGroupSeq{
				{
					newValidatorGroupSeq(testValidatorAddress, 1000, 0),
					newValidatorGroupSeq(testValidatorAddress, 1000, 0),
					newValidatorGroupSeq(testValidatorAddress, 1000, 0),
					newValidatorGroupSeq(testValidatorAddress, 1000, 0.1),
				},
			},
			expectedCount: 0,
		},
		{
			description:                "returns one missed_n_of_m system events when validator group missed 3 blocks",
			maxValidatorGroupSequences: 5,
			missedInRowThreshold:       50,
			missedForMaxThreshold:      3,
			currHeightList: []model.ValidatorGroupSeq{
				newValidatorGroupSeq(testValidatorAddress, 1000, 0),
			},
			lastForValidatorGroupList: [][]model.ValidatorGroupSeq{
				{
					newValidatorGroupSeq(testValidatorAddress, 1000, 0),
					newValidatorGroupSeq(testValidatorAddress, 1000, 0),
					newValidatorGroupSeq(testValidatorAddress, 1000, 0.1),
					newValidatorGroupSeq(testValidatorAddress, 1000, 0.1),
				},
			},
			expectedCount: 1,
			expectedKinds: []model.SystemEventKind{model.SystemEventMissedNofM},
		},
		{
			description:                "returns one missed_n_of_m system events when validator group missed 3 blocks and max < last list",
			maxValidatorGroupSequences: 3,
			missedInRowThreshold:       50,
			missedForMaxThreshold:      3,
			currHeightList: []model.ValidatorGroupSeq{
				newValidatorGroupSeq(testValidatorAddress, 1000, 0),
			},
			lastForValidatorGroupList: [][]model.ValidatorGroupSeq{
				{
					newValidatorGroupSeq(testValidatorAddress, 1000, 0),
					newValidatorGroupSeq(testValidatorAddress, 1000, 0),
					newValidatorGroupSeq(testValidatorAddress, 1000, 0.1),
					newValidatorGroupSeq(testValidatorAddress, 1000, 0.1),
					newValidatorGroupSeq(testValidatorAddress, 1000, 0.1),
				},
			},
			expectedCount: 1,
			expectedKinds: []model.SystemEventKind{model.SystemEventMissedNofM},
		},
		{
			description:                "returns no missed_n_of_m system events when count of recent not validated > maxValidatorGroupSequences",
			maxValidatorGroupSequences: 5,
			missedInRowThreshold:       50,
			missedForMaxThreshold:      3,
			currHeightList: []model.ValidatorGroupSeq{
				newValidatorGroupSeq(testValidatorAddress, 1000, 0),
			},
			lastForValidatorGroupList: [][]model.ValidatorGroupSeq{
				{
					newValidatorGroupSeq(testValidatorAddress, 1000, 0),
					newValidatorGroupSeq(testValidatorAddress, 1000, 0),
					newValidatorGroupSeq(testValidatorAddress, 1000, 0),
					newValidatorGroupSeq(testValidatorAddress, 1000, 0),
					newValidatorGroupSeq(testValidatorAddress, 1000, 0),
				},
			},
			expectedCount: 0,
		},
		{
			description:                "returns no missed_n_of_m system events when current is validated",
			maxValidatorGroupSequences: 5,
			missedInRowThreshold:       50,
			missedForMaxThreshold:      3,
			currHeightList: []model.ValidatorGroupSeq{
				newValidatorGroupSeq(testValidatorAddress, 1000, 0.1),
			},
			lastForValidatorGroupList: [][]model.ValidatorGroupSeq{
				{
					newValidatorGroupSeq(testValidatorAddress, 1000, 0),
					newValidatorGroupSeq(testValidatorAddress, 1000, 0),
					newValidatorGroupSeq(testValidatorAddress, 1000, 0),
					newValidatorGroupSeq(testValidatorAddress, 1000, 0),
					newValidatorGroupSeq(testValidatorAddress, 1000, 0),
				},
			},
			expectedCount: 0,
		},
		{
			description:                "returns error when first call to FindLastByAddress fails",
			maxValidatorGroupSequences: 3,
			missedInRowThreshold:       50,
			missedForMaxThreshold:      3,
			currHeightList: []model.ValidatorGroupSeq{
				newValidatorGroupSeq(testValidatorAddress, 1000, 0),
			},
			lastForValidatorGroupList: [][]model.ValidatorGroupSeq{
				nil,
			},
			errs:          []error{ErrCouldNotFindByAddress},
			expectedCount: 0,
			expectedErr:   ErrCouldNotFindByAddress,
		},
		{
			description:                "returns error when second call to FindLastByAddress fails",
			maxValidatorGroupSequences: 5,
			missedInRowThreshold:       3,
			missedForMaxThreshold:      5,
			currHeightList: []model.ValidatorGroupSeq{
				newValidatorGroupSeq(testValidatorAddress, 1000, 0),
				newValidatorGroupSeq("address1", 1000, 0),
			},
			lastForValidatorGroupList: [][]model.ValidatorGroupSeq{
				{
					newValidatorGroupSeq(testValidatorAddress, 1000, 0),
					newValidatorGroupSeq(testValidatorAddress, 1000, 0),
					newValidatorGroupSeq(testValidatorAddress, 1000, 0),
					newValidatorGroupSeq(testValidatorAddress, 1000, 0.1),
				},
				nil,
			},
			errs:          []error{nil, ErrCouldNotFindByAddress},
			expectedCount: 0,
			expectedErr:   ErrCouldNotFindByAddress,
		},
		{
			description:                "returns partial system events when second call to FindLastByAddress fails with ErrNotFound",
			maxValidatorGroupSequences: 3,
			missedInRowThreshold:       50,
			missedForMaxThreshold:      3,
			currHeightList: []model.ValidatorGroupSeq{
				newValidatorGroupSeq(testValidatorAddress, 1000, 0),
				newValidatorGroupSeq("address1", 1000, 0),
			},
			lastForValidatorGroupList: [][]model.ValidatorGroupSeq{
				{
					newValidatorGroupSeq(testValidatorAddress, 1000, 0),
					newValidatorGroupSeq(testValidatorAddress, 1000, 0),
					newValidatorGroupSeq(testValidatorAddress, 1000, 0.1),
					newValidatorGroupSeq(testValidatorAddress, 1000, 0.1),
				},
				nil,
			},
			errs:          []error{nil, psql.ErrNotFound},
			expectedCount: 1,
			expectedKinds: []model.SystemEventKind{model.SystemEventMissedNofM},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			validatorGroupSeqStoreMock := mock.NewMockValidatorGroupSeq(ctrl)
			accountActivitySeqStoreMock := mock.NewMockAccountActivitySeq(ctrl)

			MaxValidatorSequences = tt.maxValidatorGroupSequences
			MissedInRowThreshold = tt.missedInRowThreshold
			MissedForMaxThreshold = tt.missedForMaxThreshold

			var mockCalls []*gomock.Call
			for i, validatorGroupSeqs := range tt.lastForValidatorGroupList {
				validatorGroupSeq := tt.currHeightList[i]
				if !validatorGroupSeq.IsValidated() {
					call := validatorGroupSeqStoreMock.EXPECT().FindLastByAddress(gomock.Any(), gomock.Any())

					if len(tt.errs) >= i+1 && tt.errs[i] != nil {
						call = call.Return(nil, tt.errs[i])
					} else {
						call = call.Return(validatorGroupSeqs, nil)
					}

					mockCalls = append(mockCalls, call)
				}
			}
			gomock.InOrder(mockCalls...)

			task := NewSystemEventCreatorTask(testCfg, nil, accountActivitySeqStoreMock, validatorGroupSeqStoreMock)
			createdSystemEvents, err := task.getMissedBlocksOfValidatorGroupSequences(tt.currHeightList)
			if err == nil && tt.expectedErr != nil {
				t.Errorf("should return error")
				return
			}
			if err != nil && tt.expectedErr != err {
				t.Errorf("unexpected error, want %v; got %v", tt.expectedErr, err)
				return
			}

			if len(createdSystemEvents) != tt.expectedCount {
				t.Errorf("unexpected system event count, want %v; got %v", tt.expectedCount, len(createdSystemEvents))
				return
			}

			for i, kind := range tt.expectedKinds {
				if len(createdSystemEvents) > 0 && createdSystemEvents[i].Kind != kind {
					t.Errorf("unexpected system event kind, want %v; got %v", kind, createdSystemEvents[i].Kind)
				}
			}
		})
	}
}

func newValidatorGroupSeq(address string, membersCount int, membersAvgSigned float64) model.ValidatorGroupSeq {
	return model.ValidatorGroupSeq{
		Sequence: &model.Sequence{
			Height: testHeight,
			Time:   *types.NewTimeFromTime(time.Now()),
		},
		Address:          address,
		MembersCount:     membersCount,
		MembersAvgSigned: membersAvgSigned,
	}
}
