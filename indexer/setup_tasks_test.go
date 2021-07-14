package indexer

import (
	"context"
	"testing"

	baseClientMock "github.com/figment-networks/celo-indexer/mock/baseclient"
	figmentClientMock "github.com/figment-networks/celo-indexer/mock/client"
	"github.com/golang/mock/gomock"
)

func TestSetup_Run(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRequestCounter := baseClientMock.NewMockRequestCounter(ctrl)
	mockFigmentClient := figmentClientMock.NewMockClient(ctrl)

	task := NewSetupTask(mockFigmentClient)

	pl := &payload{CurrentHeight: 20}

	mockRequestCounter.EXPECT().InitCounter()

	mockFigmentClient.EXPECT().GetRequestCounter().Return(mockRequestCounter)

	if err := task.Run(ctx, pl); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
