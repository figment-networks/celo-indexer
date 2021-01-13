package indexer

import (
	"context"
	"github.com/figment-networks/celo-indexer/client/figmentclient"
	mock "github.com/figment-networks/celo-indexer/mock/store"
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/store/psql"
	"github.com/figment-networks/celo-indexer/types"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"reflect"
	"testing"
	"time"
)

func TestValidatorAggCreatorTask_Run(t *testing.T) {
	syncTime := types.NewTimeFromTime(time.Now())
	const syncHeight int64 = 31
	dbErr := errors.New("unexpected err")

	signed := true
	notSigned := false
	validator1 := figmentclient.Validator{Address: "acct1", Signed: nil}
	validator2 := figmentclient.Validator{Address: "acct2", Signed: &signed}
	validator3 := figmentclient.Validator{Address: "acct3", Signed: &notSigned}

	tests := []struct {
		description      string
		rawValidators    []*figmentclient.Validator
		syncable         model.Syncable
		expectErr        error
		expectValidators []model.ValidatorAgg
	}{
		{
			description:   "Adds new validator to payload.NewValidatorAggregates",
			rawValidators: []*figmentclient.Validator{&validator1, &validator2, &validator3},
			syncable: model.Syncable{
				Height: syncHeight,
				Time:   syncTime,
			},
			expectErr: nil,
			expectValidators: []model.ValidatorAgg{
				{
					Aggregate: &model.Aggregate{
						StartedAtHeight: syncHeight,
						StartedAt:       *syncTime,
						RecentAtHeight:  syncHeight,
						RecentAt:        *syncTime,
					},
					Address:                 "acct1",
					RecentAsValidatorHeight: syncHeight,
					AccumulatedUptime:       0,
					AccumulatedUptimeCount:  0,
				},
				{
					Aggregate: &model.Aggregate{
						StartedAtHeight: syncHeight,
						StartedAt:       *syncTime,
						RecentAtHeight:  syncHeight,
						RecentAt:        *syncTime,
					},
					Address:                 "acct2",
					RecentAsValidatorHeight: syncHeight,
					AccumulatedUptime:       1,
					AccumulatedUptimeCount:  1,
				},
				{
					Aggregate: &model.Aggregate{
						StartedAtHeight: syncHeight,
						StartedAt:       *syncTime,
						RecentAtHeight:  syncHeight,
						RecentAt:        *syncTime,
					},
					Address:                 "acct3",
					RecentAsValidatorHeight: syncHeight,
					AccumulatedUptime:       0,
					AccumulatedUptimeCount:  1,
				},
			},
		},
		{
			description:   "Returns err on unexpected dberr",
			rawValidators: []*figmentclient.Validator{&validator1},
			syncable: model.Syncable{
				Height: syncHeight,
				Time:   syncTime,
			},
			expectErr:        dbErr,
			expectValidators: nil,
		},
	}

	for _, tt := range tests {
		tt := tt // need to set this since running tests in parallel
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ctx := context.Background()

			dbMock := mock.NewMockValidatorAgg(ctrl)

			pld := &payload{
				RawValidators: tt.rawValidators,
				Syncable:      &tt.syncable,
			}

			for _, rawValidator := range tt.rawValidators {
				if tt.expectErr == dbErr {
					dbMock.EXPECT().FindByAddress(rawValidator.Address).Return(nil, dbErr).Times(1)
					break
				}
				dbMock.EXPECT().FindByAddress(rawValidator.Address).Return(nil, psql.ErrNotFound).Times(1)
			}

			task := NewValidatorAggCreatorTask(dbMock)
			if err := task.Run(ctx, pld); err != tt.expectErr {
				t.Errorf("unexpected error, got: %v; want: %v", err, tt.expectErr)
				return
			}

			// don't check payload if expected error
			if tt.expectErr != nil {
				return
			}

			if len(pld.NewValidatorAggregates) != len(tt.expectValidators) {
				t.Errorf("expected payload.NewValidatorAggregates to contain new accounts, got: %v; want: %v", len(pld.NewValidatorAggregates), len(tt.expectValidators))
				return
			}

			for _, expected := range tt.expectValidators {
				var found bool
				for _, got := range pld.NewValidatorAggregates {
					if got.Address == expected.Address {
						if !reflect.DeepEqual(got, expected) {
							t.Errorf("unexpected entry in payload.NewAggregatedValidators, got: %v; want: %v", got, expected)
						}
						found = true
						break
					}
				}
				if !found {
					t.Errorf("missing entry in payload.NewAggregatedValidators, want: %v", expected)
				}
			}

		})
	}

	startedAtTime := *types.NewTimeFromTime(time.Date(2020, 11, 10, 23, 0, 0, 0, time.UTC))
	const startedAtHeight int64 = 30

	updateValidatorTests := []struct {
		description      string
		rawValidators    []*figmentclient.Validator
		returnValidators []model.ValidatorAgg
		syncable         model.Syncable
		expectValidators []model.ValidatorAgg
	}{
		{
			description:   "Adds validator to payload.UpdatedValidatorAggregates",
			rawValidators: []*figmentclient.Validator{&validator1, &validator2, &validator3},
			returnValidators: []model.ValidatorAgg{
				{
					Aggregate: &model.Aggregate{
						StartedAtHeight: startedAtHeight,
						StartedAt:       startedAtTime,
						RecentAtHeight:  startedAtHeight,
						RecentAt:        startedAtTime,
					},
					Address:                 "acct1",
					RecentAsValidatorHeight: startedAtHeight,
					AccumulatedUptime:       1,
					AccumulatedUptimeCount:  1,
				},
				{
					Aggregate: &model.Aggregate{
						StartedAtHeight: startedAtHeight,
						StartedAt:       startedAtTime,
						RecentAtHeight:  startedAtHeight,
						RecentAt:        startedAtTime,
					},
					Address:                 "acct2",
					RecentAsValidatorHeight: startedAtHeight,
					AccumulatedUptime:       1,
					AccumulatedUptimeCount:  1,
				},
				{
					Aggregate: &model.Aggregate{
						StartedAtHeight: startedAtHeight,
						StartedAt:       startedAtTime,
						RecentAtHeight:  startedAtHeight,
						RecentAt:        startedAtTime,
					},
					Address:                 "acct3",
					RecentAsValidatorHeight: startedAtHeight,
					AccumulatedUptime:       1,
					AccumulatedUptimeCount:  1,
				},
			},
			syncable: model.Syncable{
				Height: syncHeight,
				Time:   syncTime,
			},
			expectValidators: []model.ValidatorAgg{
				{
					Aggregate: &model.Aggregate{
						StartedAtHeight: startedAtHeight,
						StartedAt:       startedAtTime,
						RecentAtHeight:  syncHeight,
						RecentAt:        *syncTime,
					},
					Address:                 "acct1",
					RecentAsValidatorHeight: syncHeight,
					AccumulatedUptime:       1,
					AccumulatedUptimeCount:  1,
				},
				{
					Aggregate: &model.Aggregate{
						StartedAtHeight: startedAtHeight,
						StartedAt:       startedAtTime,
						RecentAtHeight:  syncHeight,
						RecentAt:        *syncTime,
					},
					Address:                 "acct2",
					RecentAsValidatorHeight: syncHeight,
					AccumulatedUptime:       2,
					AccumulatedUptimeCount:  2,
				},
				{
					Aggregate: &model.Aggregate{
						StartedAtHeight: startedAtHeight,
						StartedAt:       startedAtTime,
						RecentAtHeight:  syncHeight,
						RecentAt:        *syncTime,
					},
					Address:                 "acct3",
					RecentAsValidatorHeight: syncHeight,
					AccumulatedUptime:       1,
					AccumulatedUptimeCount:  2,
				},
			},
		},
	}

	for _, tt := range updateValidatorTests {
		tt := tt // need to set this since running tests in parallel
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			ctx := context.Background()

			dbMock := mock.NewMockValidatorAgg(ctrl)

			pld := &payload{
				RawValidators: tt.rawValidators,
				Syncable:      &tt.syncable,
			}

			for _, validator := range tt.returnValidators {
				expect := validator
				dbMock.EXPECT().FindByAddress(validator.Address).Return(&expect, nil).Times(1)
			}

			task := NewValidatorAggCreatorTask(dbMock)
			if err := task.Run(ctx, pld); err != nil {
				t.Errorf("unexpected error, got: %v", err)
				return
			}

			if len(pld.UpdatedValidatorAggregates) != len(tt.expectValidators) {
				t.Errorf("expected payload.UpdatedValidatorAggregates to contain accounts, got: %v; want: %v", len(pld.UpdatedValidatorAggregates), len(tt.expectValidators))
				return
			}

			for _, expected := range tt.expectValidators {
				var found bool
				for _, got := range pld.UpdatedValidatorAggregates {
					if got.Address == expected.Address {
						if !reflect.DeepEqual(got, expected) {
							t.Errorf("unexpected entry in payload.UpdatedValidatorAggregates, got: %v; want: %v", got, expected)
						}
						found = true
						break
					}
				}
				if !found {
					t.Errorf("missing entry in payload.UpdatedValidatorAggregates, want: %v", expected)
				}
			}

		})
	}
}
