package main

import (
	"bufio"
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
	scanner := bufio.NewScanner(t.conn)
OUTER:
	for {
		if !scanner.Scan() {
			break OUTER
		}
		text := scanner.Text()
		_, err := t.out.Write([]byte(fmt.Sprintf("%s\n", text)))
		if err != nil {
			return fmt.Errorf("write received msg error: %w", err)
		}
	}
	log.Printf("Finished receive routine")
	return nil
}

func (t *TelnetClient) Send() error {
	scanner := bufio.NewScanner(t.in)

OUTER:
	for {
		if !scanner.Scan() {
			break OUTER
		}
		str := scanner.Text()
		_, err := t.conn.Write([]byte(fmt.Sprintf("%s\n", str)))
		if err != nil {
			return fmt.Errorf("send msg error: %w", err)
		}
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
	//	fmt.Println("conn is closed")

	if err != nil {
		return fmt.Errorf("close connection error: %w", err)
	}
	//	os.Exit(0)
	return nil
}
