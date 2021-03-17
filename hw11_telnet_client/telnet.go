package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type Client struct {
	address    string
	timeout    time.Duration
	in         io.ReadCloser
	out        io.Writer
	connection net.Conn
}

func (c *Client) Connect() error {
	var err error
	c.connection, err = net.DialTimeout("tcp", c.address, c.timeout)
	return err
}

func (c *Client) Close() error {
	err := c.connection.Close()
	return err
}

func (c *Client) Send() error {
	return c.readWrite(c.in, c.connection)
}

func (c *Client) Receive() error {
	return c.readWrite(c.connection, c.out)
}

func (c *Client) readWrite(rd io.Reader, wr io.Writer) error {
	_, err := io.Copy(wr, rd)
	return err
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}
