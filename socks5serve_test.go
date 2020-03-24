package socks5

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"testing"
	"time"
)

func TestSOCKS5_ConnectIP(t *testing.T) {
	// Create a local listener
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	go func() {
		conn, err := l.Accept()
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		defer conn.Close()

		buf := make([]byte, 4)
		if _, err := io.ReadAtLeast(conn, buf, 4); err != nil {
			t.Fatalf("err: %v", err)
		}

		if !bytes.Equal(buf, []byte("ping")) {
			t.Fatalf("bad: %v", buf)
		}
		conn.Write([]byte("pong"))
	}()
	lAddr := l.Addr().(*net.TCPAddr)

	// Create a socks server
	creds := StaticCredentials{
		"foo": "bar",
	}
	cator := FromIPUserPassAuthenticator{Credentials: creds}
	conf := &Config{
		AuthMethods: []Authenticator{cator},
		Logger:      log.New(os.Stdout, "", log.LstdFlags),
		Dial:        DialFromIP,
	}
	serv, err := New(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Start listening
	go func() {
		if err := serv.ListenAndServe("tcp", "127.0.0.1:12365"); err != nil {
			t.Fatalf("err: %v", err)
		}
	}()
	time.Sleep(10 * time.Millisecond)

	// Get a local conn
	conn, err := net.Dial("tcp", "127.0.0.1:12365")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Connect, auth and connec to local
	req := bytes.NewBuffer(nil)
	req.Write([]byte{5})
	req.Write([]byte{2, NoAuth, UserPassAuth})
	req.Write([]byte{1, 4 + 10, 'f', 'o', 'o', ':',
		'1', '0', '.', '0', '.', '0', '.', '2', '1', '2',
		3, 'b', 'a', 'r'})
	req.Write([]byte{5, 1, 0, 1, 172, 217, 10, 110})
	//req.Write([]byte{5, 1, 0, 1, 127, 0, 0, 1})

	port := []byte{8, 0}
	_ = lAddr
	//binary.BigEndian.PutUint16(port, uint16(lAddr.Port))
	req.Write(port)

	// Send a ping
	//req.Write([]byte("ping"))
	req.Write([]byte("GET / HTTP/1.1\nHost: google.com\n\n"))
	// Send all the bytes
	conn.Write(req.Bytes())

	// Verify response
	expected := []byte{
		socks5Version, UserPassAuth,
		1, authSuccess,
		5,
		0,
		0,
		1,
		172, 217, 10, 110,
		0, 0,
		'p', 'o', 'n', 'g',
	}
	//expected := [5 2 1 0 5 0 0 1 10 0 0 212 0 0 69 0 0 53 93 202 64 0 64 6 199 81 10 0 0 212 10 0 0 212 71 69 84 32 47 32 72 84 84 80 47 49 46 49 10 72]
	//out := make([]byte, len(expected))
	out := make([]byte, 50)
	fmt.Println(out)

	conn.SetDeadline(time.Now().Add(time.Second))
	if _, err := io.ReadAtLeast(conn, out, len(out)); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Ignore the port
	out[12] = 0
	out[13] = 0

	if !bytes.Equal(out, expected) {
		t.Fatalf("bad: %v", out)
	}
}
