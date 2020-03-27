package socks5

import (
	"fmt"
	"golang.org/x/net/context"
	"net"
	"strconv"
)

type key int

var ip = key(1)

func SetIP(ctx context.Context, ipVal string) context.Context {
	return context.WithValue(ctx, ip, ipVal)
}
func GetIP(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(ip).(string)
	return val, ok
}

func DialFromIP(ctx context.Context, network, addr string) (net.Conn, error) {
	fmt.Println("network:" + network)
	fmt.Println("raddr:" + addr)

	addrIP, addrPort, err := net.SplitHostPort(addr)
	if err != nil {
		panic(err)
	}
	port, err := strconv.Atoi(addrPort)
	if err != nil {
		panic(err)
	}
	raddr := net.TCPAddr{IP: net.ParseIP(addrIP), Port: port}
	ip, ok := GetIP(ctx)
	if !ok {
		fmt.Println("laddr: ASSIGNED")
		return net.DialTCP(network, nil, &raddr)
	}
	fmt.Println("laddr:" + ip)
	laddr := net.TCPAddr{IP: net.ParseIP(ip)}
	return net.DialTCP(network, &laddr, &raddr)
}
