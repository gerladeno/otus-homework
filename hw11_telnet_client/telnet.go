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
	if err != nil {
		return err
	}
	log.Printf("...Connected to %s", c.address)
	return nil
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
	if err = copyWithEOF(wr, rd); err != nil {
		if errors.Is(err, io.EOF) {
			log.Printf("...EOF")
			return nil
		}
		log.Printf("...Connection was closed by peer")
		return err
	}
	return nil
}

func copyWithEOF(dst io.Writer, src io.Reader) (err error) {
	size := 32 * 1024
	buf := make([]byte, size)
	for {
		nr, er := src.Read(buf)
		if nr > 0 { //nolint
			nw, ew := dst.Write(buf[0:nr])
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errors.New("invalid write result")
				}
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			err = er
			break
		}
	}
	return err
}
