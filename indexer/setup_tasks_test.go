package indexer

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/figment-networks/celo-indexer/client/figmentclient"
	baseClientMock "github.com/figment-networks/celo-indexer/mock/baseclient"
	figmentClientMock "github.com/figment-networks/celo-indexer/mock/client"
	"github.com/figment-networks/celo-indexer/types"
	"github.com/golang/mock/gomock"
)

func TestHeightMetaRetriever_Run(t *testing.T) {
	const chainId uint64 = 1
	const syncHeight int64 = 20

	syncTime := types.NewTimeFromTime(time.Date(2020, 11, 10, 23, 0, 0, 0, time.UTC))

	tests := []struct {
		description       string
		returnChainParams *figmentclient.ChainParams
		returnMeta        *figmentclient.HeightMeta
		result            HeightMeta
		expectErr         bool
	}{
		{
			description: "updates payload.HeightMeta",
			returnChainParams: &figmentclient.ChainParams{
				ChainId:   chainId,
				EpochSize: nil,
			},
			returnMeta: &figmentclient.HeightMeta{
				Height:      syncHeight,
				Time:        uint64(syncTime.Time.UTC().Unix()),
				Epoch:       nil,
				LastInEpoch: nil,
			},
			result: HeightMeta{
				ChainId:     chainId,
				Height:      syncHeight,
				Time:        syncTime,
				Epoch:       nil,
				EpochSize:   nil,
				LastInEpoch: nil,
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRequestCounter := baseClientMock.NewMockRequestCounter(ctrl)
			mockFigmentClient := figmentClientMock.NewMockClient(ctrl)

			task := NewHeightMetaRetrieverTask(mockFigmentClient)

			pl := &payload{CurrentHeight: 20}

			mockRequestCounter.EXPECT().InitCounter()

			mockFigmentClient.EXPECT().GetRequestCounter().Return(mockRequestCounter)
			mockFigmentClient.EXPECT().GetChainParams(ctx).Return(tt.returnChainParams, nil).Times(1)
			mockFigmentClient.EXPECT().GetMetaByHeight(ctx, pl.CurrentHeight).Return(tt.returnMeta, nil).Times(1)

			if err := task.Run(ctx, pl); err != nil && !tt.expectErr {
				t.Errorf("unexpected error, want %v; got %v", tt.expectErr, err)
				return
			}

			if !reflect.DeepEqual(pl.HeightMeta, tt.result) {
				t.Errorf("want: %+v, got: %+v", tt.result, pl.HeightMeta)
				return
			}
		})
	}
}
