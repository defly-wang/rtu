package config

import "net"

type ProxyConnect struct {
	Dtumac       string
	DtuIp        string
	Targetport   uint
	DtuConnected bool
	ProxyIp      string
	ProxyConn    net.Conn
	DirectPort   uint
	DirectConn   net.Conn
}

type Proxy struct {
	Remoteport int //22
	Localport  int
	Proxyport  int
	Ip         string
}
