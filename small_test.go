package socks5

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

func Test5_Connect(t *testing.T) {
	// Get a local conn
	conn, err := net.Dial("tcp", "185.243.217.108:9998")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer conn.Close()

	// Connect, auth and connec to local
	req := bytes.NewBuffer(nil)
	req.Write([]byte{5})

	// Send a ping
	req.Write([]byte("ping"))

	// Send all the bytes
	conn.Write(req.Bytes())

	out := make([]byte, 20)

	conn.SetDeadline(time.Now().Add(time.Second))
	if _, err := io.ReadAtLeast(conn, out, 4); err != nil {
		t.Fatalf("err: %v", err)
	}
	fmt.Println(string(out))
}
