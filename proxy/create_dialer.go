package proxy

import (
	"fmt"
	"golang.org/x/net/proxy"
	"net"
)

var (
	proxyDialer proxy.Dialer
	netDialer   net.Dialer
)

// ConfigureDialer sets up a proxy dialer if the first argument is an address
// and a net.Dialer if the first argument is ""
func ConfigureDialer(proxyAddress string) (err error) {
	if proxyAddress != "" {
		fmt.Println("Using proxy")
		proxyDialer, err = proxy.SOCKS5("tcp", proxyAddress, nil, proxy.Direct)
	} else {
		netDialer = net.Dialer{}
	}

	return
}

// GetDialer creates a single interface to get a dialer
func GetDialer(address string) (net.Conn, error) {
	if proxyDialer != nil {
		return proxyDialer.Dial("tcp", address)
	} else {
		return netDialer.Dial("tcp", address)
	}
}
