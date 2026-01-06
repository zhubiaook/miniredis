package client

import "net"

type Options struct {
	Addr string
	DB   int
}

type Client struct {
	conn net.Conn
}

func NewClient(opts Options) (*Client, error) {
	conn, err := net.Dial("tcp", opts.Addr)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Do(args ...string) ([]byte, error) {
	req := Encode(args)
	if _, err := c.conn.Write(req); err != nil {
		return nil, err
	}

	resp := make([]byte, 1024)
	if _, err := c.conn.Read(resp); err != nil {
		return nil, err
	}
	return resp, nil
}
