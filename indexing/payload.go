package indexing

import (
	"sync"

	"github.com/figment-networks/indexing-engine/pipeline"
)

var (
	payloadPool = sync.Pool{
		New: func() interface{} {
			return new(payload)
		},
	}

	_ pipeline.PayloadFactory = (*payloadFactory)(nil)
	_ pipeline.Payload        = (*payload)(nil)
)

func NewPayloadFactory() *payloadFactory {
	return &payloadFactory{}
}

type payloadFactory struct{}

func (*payloadFactory) GetPayload() pipeline.Payload {
	return payloadPool.Get().(*payload)
}

type payload struct {
	currentHeight int64
}

func (p *payload) SetCurrentHeight(height int64) {
	p.currentHeight = height
}

func (p *payload) GetCurrentHeight() int64 {
	return p.currentHeight
}

func (p *payload) MarkAsProcessed() {
	payloadPool.Put(p)
}
