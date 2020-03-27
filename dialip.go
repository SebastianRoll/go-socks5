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
		return net.Dial(network, addr)
	}
	ip = ip
	fmt.Println("network:" + network)
	fmt.Println("laddr:" + ip)
	fmt.Println("raddr:" + addr)
	var (
		laddr net.IPAddr = net.IPAddr{IP: net.ParseIP(ip)}
		raddr net.IPAddr = net.IPAddr{IP: net.ParseIP(addr)}
	)
	_ = laddr
	return net.DialIP(network, &laddr, &raddr)
}
