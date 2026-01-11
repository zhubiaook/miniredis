package main

import (
	"fmt"
	"net"

	"github.com/zhubiaook/miniredis/pkg/encoding"
)

type Options struct {
	// Addr is the address formated as host:port
	Addr string

	// Username is used to authenticate the current connection
	// with one of the connections defined in the ACL list when connecting
	// to a Redis 6.0 instance, or greater, that is using the Redis ACL system.
	Username string

	// Password is an optional password. Must match the password specified in the
	// `requirepass` server configuration option (if connecting to a Redis 5.0 instance, or lower),
	// or the User Password when connecting to a Redis 6.0 instance, or greater,
	// that is using the Redis ACL system.
	Password string

	// DB is the database to be selected after connecting to the server.
	DB int
}

type Client struct {
	conn net.Conn
	opt  Options
}

func NewClient(opt Options) (*Client, error) {
	conn, err := net.Dial("tcp", opt.Addr)
	if err != nil {
		return nil, fmt.Errorf("connect to server failed: %w", err)
	}

	return &Client{
		conn: conn,
		opt:  opt,
	}, nil
}

func (c *Client) Do(args []string) ([]string, error) {
	if err := encoding.EncodeWrite(c.conn, args); err != nil {
		return nil, fmt.Errorf("encode write failed: %w", err)
	}

	var out []string
	if err := encoding.DecodeRead(c.conn, &out); err != nil {
		return nil, fmt.Errorf("decode read failed: %w", err)
	}

	return out, nil
}

func (c *Client) Close() {
	c.conn.Close()
}
