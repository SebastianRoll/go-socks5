package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/SebastianRoll/go-socks5"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

type handler func(w http.ResponseWriter, r *http.Request)

func basicAuth(pass handler) handler {

	return func(w http.ResponseWriter, r *http.Request) {

		auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 || !validate(pair[0], pair[1]) {
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}

		pass(w, r)
	}
}

func validate(username, password string) bool {
	fmt.Println(username)
	fmt.Println(os.Getenv("SOCKS5_USER"))
	if username == os.Getenv("SOCKS5_USER") && password == os.Getenv("SOCKS5_PASSWORD") {
		return true
	}
	return false
}

type InterfaceResponse struct {
	Name  string
	Flags []string
	Addrs []string
}

func interfaces(w http.ResponseWriter, req *http.Request) {
	fmt.Println("In /interfaces")
	w.Header().Set("Content-Type", "application/json")

	ints, err := net.Interfaces()
	response := []InterfaceResponse{}
	if err != nil {
		panic(err)
	}
	for _, s := range ints {
		intresp := InterfaceResponse{}
		intresp.Name = s.Name

		if s.Flags&net.FlagUp != 0 {
			intresp.Flags = append(intresp.Flags, "FlagUp")
		}
		if s.Flags&net.FlagBroadcast != 0 {
			intresp.Flags = append(intresp.Flags, "FlagBroadcast")
		}
		if s.Flags&net.FlagLoopback != 0 {
			intresp.Flags = append(intresp.Flags, "FlagLoopback")
		}
		if s.Flags&net.FlagPointToPoint != 0 {
			intresp.Flags = append(intresp.Flags, "FlagPointToPoint")
		}
		if s.Flags&net.FlagMulticast != 0 {
			intresp.Flags = append(intresp.Flags, "FlagMulticast")
		}

		addrs, err := s.Addrs()
		if err != nil {
			panic(err)
		}
		adds := []string{}
		for _, a := range addrs {
			adds = append(adds, a.String())
		}
		intresp.Addrs = adds
		response = append(response, intresp)
	}
	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(js)

}

func main() {
	// Create a local listener
	l, err := net.Listen("tcp", ":9998")
	if err != nil {
		panic(err)
	}
	fmt.Println("PONG server ja")
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				panic(err)
			}
			fmt.Println("in PONG server")

			buf := make([]byte, 5)
			if _, err := io.ReadAtLeast(conn, buf, 4); err != nil {
				panic(err)
			}
			fmt.Printf(string(buf))

			//if !bytes.Equal(buf, []byte("ping")) {
			//	t.Fatalf("bad: %v", buf)
			//}
			conn.Write([]byte("pong"))
			conn.Close()

		}
	}()
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
	go func() {
		if err := serv.ListenAndServe("tcp", ":8989"); err != nil {
			panic(err)
		}
	}()

	//http.HandleFunc("/interfaces", basicAuth(interfaces))
	http.HandleFunc("/interfaces", interfaces)
	http.ListenAndServe(":8998", nil)

}
