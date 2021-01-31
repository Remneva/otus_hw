package main

import (
	"fmt"
	"io"
	"log"
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
	_, err := io.Copy(t.out, t.conn)
	if err != nil {
		return fmt.Errorf("write received msg error: %w", err)
	}
	log.Printf("Finished receive routine")
	return nil
}

func (t *TelnetClient) Send() error {
	_, err := io.Copy(t.conn, t.in)
	if err != nil {
		return fmt.Errorf("send msg error: %w", err)
	}
	log.Printf("Finished send routine")
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
	err := t.conn.Close()
	if err != nil {
		return fmt.Errorf("close connection error: %w", err)
	}
	return nil
}
