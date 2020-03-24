package main

import (
	"github.com/sebastianroll/go-socks5"
	"log"
	"os"
)

func main() {

	// Create a socks server
	creds := socks5.StaticCredentials{
		os.Getenv("SOCKS5_USER"): os.Getenv("SOCKS5_PASSWORD"),
	}
	cator := socks5.FromIPUserPassAuthenticator{Credentials: creds}
	conf := &socks5.Config{
		AuthMethods: []socks5.Authenticator{cator},
		Logger:      log.New(os.Stdout, "", log.LstdFlags),
		Dial:        socks5.DialFromIP,
	}
	serv, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	// Start listening
	if err := serv.ListenAndServe("tcp", "127.0.0.1:8989"); err != nil {
		panic(err)
	}

}
