package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type TelnetClient struct {
	address  string
	timeout  time.Duration
	conn     net.Conn
	listener net.Listener
	ctx      context.Context
	in       io.ReadCloser
	out      io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	ctx := context.Background()
	t := TelnetClient{
		address: address,
		timeout: timeout,
		ctx:     ctx,
		in:      in,
		out:     out,
	}

	return t
}

func (t *TelnetClient) Receive() error {
	scanner := bufio.NewScanner(t.conn)
OUTER:
	for {
		select {
		case <-t.ctx.Done():
			break OUTER
		default:
			if !scanner.Scan() {
				log.Printf("CANNOT SCAN")
				break OUTER
			}
			text := scanner.Text()
			fmt.Println(text)
			_, err := t.out.Write([]byte(fmt.Sprintf("%s\n", text)))
			if err != nil {
				return fmt.Errorf("write received msg error: %w", err)
			}
		}
	}
	log.Printf("Finished receive routine")
	return nil
}

func (t *TelnetClient) Send() error {
	scanner := bufio.NewScanner(t.in)
OUTER:
	for {
		select {
		case <-t.ctx.Done():
			break OUTER
		default:
			if !scanner.Scan() {
				break OUTER
			}
			str := scanner.Text()
			fmt.Println(str)
			_, err := t.conn.Write([]byte(fmt.Sprintf("%s\n", str)))
			if err != nil {
				return fmt.Errorf("send msg error: %w", err)
			}
		}
	}
	log.Printf("Finished send routine")
	return nil
}

func (t *TelnetClient) Connect() error {
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		return fmt.Errorf("announces on the local network address error: %w", err)
	}
	t.listener = listener
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
	err = t.listener.Close()
	if err != nil {
		return fmt.Errorf("close listener error: %w", err)
	}
	return nil
}
