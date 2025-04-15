package common

import (
	"net"
	"time"
)

func TestTCPWithTimeout(protocol, address string, timeout time.Duration) (net.Conn, error) {
	dial := &net.Dialer{Timeout: timeout}
	return WrapTcp(protocol, address, dial)
}

func WrapTcp(protocol, address string, dial *net.Dialer) (net.Conn, error) {
	//var conn net.Conn
	// todo: add proxy: like socks5
	conn, err := dial.Dial(protocol, address)
	if err != nil {
		return conn, err
	}
	return conn, nil
}
