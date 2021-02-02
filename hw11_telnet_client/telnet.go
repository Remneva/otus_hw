package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient struct {
	address string
	timeout time.Duration
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	t := TelnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
	return t
}

func (t *TelnetClient) Receive() error {
	if _, err := io.Copy(t.out, t.conn); err != nil {
		return fmt.Errorf("write received msg error: %w", err)
	}
	return nil
}

func (t *TelnetClient) Send() error {
	if _, err := io.Copy(t.conn, t.in); err != nil {
		return fmt.Errorf("send msg error: %w", err)
	}
	return nil
}

func (t *TelnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return fmt.Errorf("dial connection error: %w", err)
	}
	t.conn = conn
	return nil
}

func (t *TelnetClient) Close() error {
	if err := t.conn.Close(); err != nil {
		return fmt.Errorf("close connection error: %w", err)
	}
	return nil
}
