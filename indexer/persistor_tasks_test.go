package indexer

import (
	"context"
	"fmt"
	"testing"
	"time"

	mock "github.com/figment-networks/celo-indexer/mock/store"
	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/types"
	"github.com/golang/mock/gomock"
)

func TestSyncerPersistor_Run(t *testing.T) {
	sync := &model.Syncable{
		Height: 20,
		Time:   types.NewTimeFromTime(time.Now()),
	}
	t.Parallel()

	tests := []struct {
		description string
		expectErr   error
	}{
		{"calls db with syncable", nil},
		{"returns error if database errors", fmt.Errorf("test err")},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			dbMock := mock.NewMockSyncables(ctrl)

			task := NewSyncerPersistorTask(dbMock)

			pl := &payload{
				CurrentHeight: 20,
				Syncable:      sync,
			}

			dbMock.EXPECT().CreateOrUpdate(sync).Return(tt.expectErr).Times(1)

			if err := task.Run(ctx, pl); err != tt.expectErr {
				t.Errorf("want %v; got %v", tt.expectErr, err)
			}
		})
	}
}

func TestBlockSeqPersistor_Run(t *testing.T) {
	seq := &model.BlockSeq{
		Sequence: &model.Sequence{
			Height: 20,
			Time:   *types.NewTimeFromTime(time.Date(1987, 12, 11, 14, 0, 0, 0, time.UTC)),
		},
		TxCount: 10,
	}

	tests := []struct {
		description string
		expectErr   error
	}{
		{"calls db with block sequence", nil},
		{"returns error if database errors", fmt.Errorf("test err")},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("[new] %v", tt.description), func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			dbMock := mock.NewMockBlockSeq(ctrl)

			task := NewBlockSeqPersistorTask(dbMock)

			pl := &payload{
				CurrentHeight:    20,
				NewBlockSequence: seq,
			}

			dbMock.EXPECT().Create(seq).Return(tt.expectErr).Times(1)

			if err := task.Run(ctx, pl); err != tt.expectErr {
				t.Errorf("want %v; got %v", tt.expectErr, err)
			}
		})

		t.Run(fmt.Sprintf("[updated] %v", tt.description), func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			dbMock := mock.NewMockBlockSeq(ctrl)

			task := NewBlockSeqPersistorTask(dbMock)

			pl := &payload{
				CurrentHeight:        20,
				UpdatedBlockSequence: seq,
			}

			dbMock.EXPECT().Save(seq).Return(tt.expectErr).Times(1)

			if err := task.Run(ctx, pl); err != tt.expectErr {
				t.Errorf("want %v; got %v", tt.expectErr, err)
			}
		})
	}
}

func TestValidatorSeqPersistor_Run(t *testing.T) {
	signed := true
	notSigned := false
	seqs := []model.ValidatorSeq{
		{Sequence: &model.Sequence{Height: 20}, Address: "acct1", Signed: &notSigned},
		{Sequence: &model.Sequence{Height: 20}, Address: "acct2", Signed: &signed},
		{Sequence: &model.Sequence{Height: 20}, Address: "acct3", Signed: &notSigned},
	}

	tests := []struct {
		description string
		expectErr   error
	}{
		{"calls db with all validator sequences", nil},
		{"returns error if database errors", fmt.Errorf("db err")},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			dbMock := mock.NewMockValidatorSeq(ctrl)

			task := NewValidatorSeqPersistorTask(dbMock)

			pl := &payload{
				Syncable:           &model.Syncable{Height: 20},
				ValidatorSequences: seqs,
			}

			dbMock.EXPECT().BulkUpsert(seqs).Return(tt.expectErr).Times(1)

			if err := task.Run(ctx, pl); err != tt.expectErr {
				t.Errorf("want %v; got %v", tt.expectErr, err)
			}
		})
	}
}

func TestValidatorGroupSeqPersistor_Run(t *testing.T) {
	seqs := []model.ValidatorGroupSeq{
		{Sequence: &model.Sequence{Height: 20}, Address: "acct1", Name: "test1"},
		{Sequence: &model.Sequence{Height: 20}, Address: "acct2", Name: "test2"},
		{Sequence: &model.Sequence{Height: 20}, Address: "acct3", Name: "test3"},
	}

	tests := []struct {
		description string
		expectErr   error
	}{
		{"calls db with all validator sequences", nil},
		{"returns error if database errors", fmt.Errorf("db err")},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			dbMock := mock.NewMockValidatorGroupSeq(ctrl)

			task := NewValidatorGroupSeqPersistorTask(dbMock)

			pl := &payload{
				Syncable:                &model.Syncable{Height: 20},
				ValidatorGroupSequences: seqs,
			}

			dbMock.EXPECT().BulkUpsert(seqs).Return(tt.expectErr).Times(1)

			if err := task.Run(ctx, pl); err != tt.expectErr {
				t.Errorf("want %v; got %v", tt.expectErr, err)
			}
		})
	}
}
