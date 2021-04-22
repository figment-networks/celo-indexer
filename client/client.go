package client

type Client interface {
	GetName() string
	Close()
}

type RequestCounter interface {
	InitCounter()
	IncrementCounter()
	GetCounter() uint64
}
