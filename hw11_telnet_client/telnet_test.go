package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})

	t.Run("connect to wrong port", func(t *testing.T) {
		d := 10 * time.Second
		client := NewTelnetClient("localhost:121231321", d, os.Stdin, os.Stdout)
		require.Error(t, client.Connect())
	})

	t.Run("EOF", func(t *testing.T) {
		s, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, s.Close())
		}()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout := 10 * time.Second

			c := NewTelnetClient(s.Addr().String(), timeout, ioutil.NopCloser(in), out)
			require.NoError(t, c.Connect())
			defer func() {
				require.NoError(t, c.Close())
			}()

			in.WriteString("test\n")
			err = c.Send()
			require.NoError(t, err)
		}()

		go func() {
			defer wg.Done()

			c, err := s.Accept()
			require.NoError(t, err)
			require.NotNil(t, c)
			defer func() {
				require.NoError(t, c.Close())
			}()

			b := make([]byte, 1024)
			_, err = c.Read(b)
			require.NoError(t, err)

			_, err = c.Read(b)
			require.EqualError(t, err, io.EOF.Error())
		}()

		wg.Wait()
	})
}
