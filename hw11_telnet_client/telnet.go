package main

import (
	"errors"
	"io"
	"log"
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

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (c *Client) Connect() error {
	var err error
	c.connection, err = net.DialTimeout("tcp", c.address, c.timeout)
	log.Printf("...Connected to %s", c.address)
	return err
}

func (c *Client) Close() (err error) {
	if err = c.in.Close(); err != nil {
		return
	}
	err = c.connection.Close()
	return
}

func (c *Client) Send() (err error) {
	err = c.readWrite(c.in, c.connection)
	return
}

func (c *Client) Receive() (err error) {
	err = c.readWrite(c.connection, c.out)
	return
}

func (c *Client) readWrite(rd io.Reader, wr io.Writer) (err error) {
	if c.connection == nil {
		return
	}
	if _, err := io.Copy(wr, rd); err != nil {
		if errors.Is(err, io.EOF) {
			log.Printf("...EOF")
		}
		log.Printf("...Connection was closed by peer")
		return err
	}
	return nil
}
