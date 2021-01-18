package indexer

import (
	"context"
	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/config"
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/types"
	"math/big"
	"reflect"
	"testing"
	"time"
)

func TestValidatorSeqCreator_Run(t *testing.T) {
	const syncHeight int64 = 20
	testCfg := &config.Config{
		FirstBlockHeight: 1,
	}

	syncTime := *types.NewTimeFromTime(time.Date(2020, 11, 10, 23, 0, 0, 0, time.UTC))

	seq := &model.Sequence{
		Height: syncHeight,
		Time:   syncTime,
	}

	tests := []struct {
		description string
		raw         []*figmentclient.Validator
		expect      []model.ValidatorSeq
		expectErr   bool
	}{
		{
			description: "updates payload.ValidatorSequences",
			raw: []*figmentclient.Validator{
				{Address: "validator1", Score: big.NewInt(100), Signed: nil},
			},
			expect: []model.ValidatorSeq{
				{
					Sequence: seq,
					Address:  "validator1",
					Score:    types.NewQuantityFromInt64(100),
					Signed:   nil,
				},
			},

			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctx := context.Background()

			task := NewValidatorSeqCreatorTask(testCfg)

			pl := &payload{
				CurrentHeight: syncHeight,
				Syncable: &model.Syncable{
					Height: syncHeight,
					Time:   &syncTime,
				},
				RawValidators: tt.raw,
			}

			if err := task.Run(ctx, pl); err != nil && !tt.expectErr {
				t.Errorf("unexpected error, want %v; got %v", tt.expectErr, err)
				return
			}

			// skip payload check if there's an error
			if tt.expectErr {
				return
			}

			if len(pl.ValidatorSequences) != (len(tt.raw)) {
				t.Errorf("expected payload.ValidatorSequences to contain all validator seqs, got: %v; want: %v", len(pl.ValidatorSequences), len(tt.raw))
				return
			}

			for _, expectVal := range tt.expect {
				var found bool
				for _, val := range pl.ValidatorSequences {
					if val.Address == expectVal.Address {
						if !reflect.DeepEqual(val, expectVal) {
							t.Errorf("unexpected entry in payload.ValidatorSequences, got: %v; want: %v", val, expectVal)
						}
						found = true
						break
					}
				}
				if !found {
					t.Errorf("missing entry in payload.ValidatorSequences, want: %v", expectVal)
				}
			}
		})
	}
}

func TestValidatorGroupSeqCreator_Run(t *testing.T) {
	const syncHeight int64 = 20
	testCfg := &config.Config{
		FirstBlockHeight: 1,
	}

	syncTime := *types.NewTimeFromTime(time.Date(2020, 11, 10, 23, 0, 0, 0, time.UTC))

	seq := &model.Sequence{
		Height: syncHeight,
		Time:   syncTime,
	}

	tests := []struct {
		description string
		raw         []*figmentclient.ValidatorGroup
		expect      []model.ValidatorGroupSeq
		expectErr   bool
	}{
		{
			description: "updates payload.ValidatorSequences",
			raw: []*figmentclient.ValidatorGroup{
				{
					Address:          "validator1",
					Commission:       big.NewInt(100),
					VotingCap:        big.NewInt(101),
					PendingVotes:     big.NewInt(102),
					ActiveVotes:      big.NewInt(103),
				},
			},
			expect: []model.ValidatorGroupSeq{
				{
					Sequence: seq,
					Address:  "validator1",
					Commission: types.NewQuantityFromInt64(100),
					VotingCap: types.NewQuantityFromInt64(101),
					PendingVotes: types.NewQuantityFromInt64(102),
					ActiveVotes: types.NewQuantityFromInt64(103),
				},
			},

			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctx := context.Background()

			task := NewValidatorGroupSeqCreatorTask(testCfg)

			pl := &payload{
				CurrentHeight: syncHeight,
				Syncable: &model.Syncable{
					Height: syncHeight,
					Time:   &syncTime,
				},
				RawValidatorGroups: tt.raw,
			}

			if err := task.Run(ctx, pl); err != nil && !tt.expectErr {
				t.Errorf("unexpected error, want %v; got %v", tt.expectErr, err)
				return
			}

			// skip payload check if there's an error
			if tt.expectErr {
				return
			}

			if len(pl.ValidatorGroupSequences) != (len(tt.raw)) {
				t.Errorf("expected payload.ValidatorGroupSequences to contain all validator seqs, got: %v; want: %v", len(pl.ValidatorGroupSequences), len(tt.raw))
				return
			}

			for _, expectVal := range tt.expect {
				var found bool
				for _, val := range pl.ValidatorGroupSequences {
					if val.Address == expectVal.Address {
						if !reflect.DeepEqual(val, expectVal) {
							t.Errorf("unexpected entry in payload.ValidatorGroupSequences, got: %v; want: %v", val, expectVal)
						}
						found = true
						break
					}
				}
				if !found {
					t.Errorf("missing entry in payload.ValidatorGroupSequences, want: %v", expectVal)
				}
			}
		})
	}
}
