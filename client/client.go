package client

type Client interface {
	GetName() string
	Close()
}