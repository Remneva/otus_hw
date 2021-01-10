package main

import (
	"bytes"
	"io/ioutil"
	"net"
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
	t.Run("connection error", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		timeout, err := time.ParseDuration("5s")
		require.NoError(t, err)

		client := NewTelnetClient("ololo.com:80", timeout, ioutil.NopCloser(in), out)
		err = client.Connect()
		require.Error(t, err)
	})

	t.Run("close connection error", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		timeout, err := time.ParseDuration("0s")
		require.NoError(t, err)

		client := NewTelnetClient("rbc.ru:80", timeout, ioutil.NopCloser(in), out)
		require.NoError(t, client.Connect())

		err = client.Close()
		require.NoError(t, err)
		err = client.Close()
		require.Error(t, err)
	})

	//	t.Run("close connection error", func(t *testing.T) {
	//		l, err := net.Listen("tcp", "127.0.0.1:")
	//		require.NoError(t, err)
	//		defer func() { require.NoError(t, l.Close()) }()
	//
	//		timeout, err := time.ParseDuration("0s")
	//		require.NoError(t, err)
	//		i := &bytes.Buffer{}
	//		in := ioutil.NopCloser(i)
	//		client := NewTelnetClient("rbc.ru:80", timeout, ioutil.NopCloser(in), os.Stdout)
	//
	//		require.Equal(t, nil, client.conn)
	//		require.NoError(t, client.Connect())
	//		i.WriteString("hello\n")
	//		os.Stdin.Close()
	//	//	err = client.Close()
	//		fmt.Println(err)
	//		err = in.Close()
	//		close(read(os.Stdin))
	//		require.NoError(t, err)
	////		time.Sleep(5 * time.Second)
	//	//	require.Equal(t, nil, client.conn)
	//
	//		require.Equal(t, nil, client.conn)
	////		err = client.Send()
	//		require.Error(t, err)
	//
	//	})
}
