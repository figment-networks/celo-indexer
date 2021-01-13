package indexer

import (
	"context"
	"errors"
	"github.com/figment-networks/celo-indexer/client/figmentclient"
	mock "github.com/figment-networks/celo-indexer/mock/client"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func TestBlockFetcher_Run(t *testing.T) {
	tests := []struct {
		description string
		returnBlock *figmentclient.Block
		result      error
	}{
		{"returns error if client errors", nil, errors.New("test error")},
		{"updates payload.RawBlock", &figmentclient.Block{Height: 20}, nil},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			mockClient := mock.NewMockClient(ctrl)
			task := NewBlockFetcherTask(mockClient)

			pl := &payload{CurrentHeight: 20}

			mockClient.EXPECT().GetBlockByHeight(ctx, pl.CurrentHeight).Return(tt.returnBlock, tt.result).Times(1)

			if result := task.Run(ctx, pl); result != tt.result {
				t.Errorf("want %v; got %v", tt.result, result)
				return
			}

			// skip payload check if there's an error
			if tt.result != nil {
				return
			}

			if !reflect.DeepEqual(pl.RawBlock, tt.returnBlock) {
				t.Errorf("want: %+v, got: %+v", tt.returnBlock, pl.RawBlock)
				return
			}
		})
	}
}
