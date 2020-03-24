package socks5

import (
	"fmt"
	"golang.org/x/net/context"
	"net"
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
	ip, ok := GetIP(ctx)
	if !ok {
		panic("IP was not found")
	}
	fmt.Println(ip)
	fmt.Println(addr)
	var (
		laddr net.IPAddr = net.IPAddr{IP: net.ParseIP(ip)}
		raddr net.IPAddr = net.IPAddr{IP: net.ParseIP(addr)}
	)
	return net.DialIP(network, &laddr, &raddr)
}
