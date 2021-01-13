package client

type Client interface {
	GetName() string
	Close()
}

type RequestCounter interface {
	InitCounter()
	IncrementCounter() uint64
	GetCounter() uint64
}