package client

import (
	"context"

	"google.golang.org/grpc"

	"github.com/figment-networks/celo-indexer/grpc/validator"
)

var (
	_ ValidatorClient = (*validatorClient)(nil)
)

type ValidatorClient interface {
	GetByHeight(int64) (*validator.GetByHeightResponse, error)
}

func NewValidatorClient(conn *grpc.ClientConn) ValidatorClient {
	return &validatorClient{
		client: validator.NewValidatorServiceClient(conn),
	}
}

type validatorClient struct {
	client validator.ValidatorServiceClient
}

func (r *validatorClient) GetByHeight(h int64) (*validator.GetByHeightResponse, error) {
	ctx := context.Background()

	return r.client.GetByHeight(ctx, &validator.GetByHeightRequest{Height: h})
}
