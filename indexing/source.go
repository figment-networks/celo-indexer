package indexing

import (
	"context"

	"github.com/figment-networks/indexing-engine/pipeline"
)

func NewSource() pipeline.Source {
	return &source{}
}

type source struct {
	err error
}

func (s *source) Next(ctx context.Context, p pipeline.Payload) bool {
	return false
}

func (s *source) Current() int64 {
	return 0
}

func (s *source) Err() error {
	return s.err
}
