package client

import (
	"google.golang.org/grpc"
)

func New(connStr string) (*Client, error) {
	conn, err := grpc.Dial(connStr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,

		Validator: NewValidatorClient(conn),
	}, nil
}

type Client struct {
	conn *grpc.ClientConn

	Validator ValidatorClient
}

func (c *Client) Close() error {
	return c.conn.Close()
}
