package main

import (
	"config"
	"control-server/common"
	"fmt"
	"net"
	"time"
)

func LocalServerListen(port int) *net.TCPListener {
	listenserver, _ := net.ResolveTCPAddr("tcp4", fmt.Sprint(":", port))
	listener, err := net.ListenTCP("tcp", listenserver)
	if err == nil {
		return listener
	} else {
		return nil
	}
}

func LocalServer(port int) {
	//1.port（例如6001） 本地端口，用于监听本地接入，打开远程透传时开启并监听
	//2.连接接入后（本地电脑远程访问或控制），开启代理端口（例如6201），同时打开6001-6201透明代理。通知rtuclient可以接入
	//3.rtuclient接到通知后，一面连接远程服务器（例如22），一面接入代理（例如6201），同时打开22-6201透明代理。返回接入成功

	listenserver, _ := net.ResolveTCPAddr("tcp4", fmt.Sprint(":", port))
	listener, err := net.ListenTCP("tcp", listenserver)
	if err != nil {
		return
	}

	defer listener.Close()

	//存listener
	setRmtConnectListener(port, listener)

	for {
		conn, err := listener.Accept()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "use of closed network connection" {
				// 这是一个已知的关闭错误，可以安全忽略
				//log.Println("Listen socket is closed.")
				break
			}
			//var errr  error
			//fmt.Println("accept:", err)
			//删除该连接
			continue
		}

		freeproxyport := GetFreeProxyPort()

		proxylistener := ProxyServerListen(freeproxyport)
		if proxylistener == nil {
			conn.Close()
			continue
		}

		setProxyListen(port, freeproxyport, proxylistener)
		go ProxyServer(port, freeproxyport, proxylistener)

		//此处可优化为进程等待，等待ProxyServer 监听成功后执行
		//此等待很有效
		time.Sleep(250 * time.Millisecond)

		deccon := findRmtClientFromLocalport(port)
		cmd := config.Cmd{
			Cmd: config.CMD_Remote,
			Dec: deccon.Ip,
			Data: config.Proxy{
				Localport:  port,
				Proxyport:  freeproxyport,
				Remoteport: deccon.Remoteport,
			},
		}

		id := findRtuClient(cmd)

		if id == 0 {
			conn.Close()
			continue
		}

		if !SendDataJson(getRtuClient(id).Conn, cmd) {
			conn.Close()
			continue
		}

		//存储 连接、proxy端口6201
		addProxyRmtClient(port, RmtClient{
			Localconnert: conn,
			Proxyport:    freeproxyport,
		})

	}
}

func ProxyServerListen(proxyport int) *net.TCPListener {
	proxyserver, _ := net.ResolveTCPAddr("tcp4", fmt.Sprint(":", proxyport))
	proxylistener, err := net.ListenTCP("tcp", proxyserver)
	if err == nil {
		return proxylistener
	} else {
		return nil
	}
}

func ProxyServer(localport int, proxyport int, proxylistener *net.TCPListener) {
	/*
		proxyserver, _ := net.ResolveTCPAddr("tcp4", fmt.Sprint(":", proxyport))
		proxylistener, err := net.ListenTCP("tcp", proxyserver)
		if err != nil {
			return
		}
	*/
	defer proxylistener.Close()

	//setProxyListen(localport, proxyport, proxylistener)

	//需要循环吗？
	for {
		conn, err := proxylistener.Accept()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "use of closed network connection" {
				// 这是一个已知的关闭错误，可以安全忽略
				//log.Println("Listen socket is closed.")
				break
			}
			fmt.Println("accept:", err)
			//删除该连接
			continue
		}

		//save
		//setProxyConect(localport, proxyport, conn)
		//fmt.Println("port:", localport, "pport:", proxyport)
		localconn := findRmtLocalConnect(localport, proxyport)
		if localconn != nil {
			setProxyConect(localport, proxyport, conn)
			go copyData(conn, localconn, localport, proxyport)
			go copyData(localconn, conn, localport, proxyport)
			//go haneleProxyServer(localport, proxyport, conn, localconn)
		}

		//time.Sleep(10 * time.Millisecond)

	}
}

/*
func haneleProxyServer(localport, proxyport int, conn, localconn net.Conn) {

	//localconn := findRmtLocalConnect(localport, proxyport)

	//if localconn != nil {
	go copyData(conn, localconn, localport, proxyport)
	go copyData(localconn, conn, localport, proxyport)

	//}

}
*/
func copyData(source, dest net.Conn, localport, proxyport int) {
	buf := make([]byte, 4096)
	for {
		if source == nil || dest == nil {
			return
		}
		n, err := source.Read(buf)
		if err != nil {
			break
		}

		_, err = dest.Write(buf[:n])
		if err != nil {
			break
		}
	}
	defer func() {
		delRmtClient(localport, proxyport)
		if source != nil {
			source.Close()
			source = nil
		}
		if dest != nil {
			dest.Close()
			dest = nil
		}

	}()
}

func GetFreeLocalPort() int {

	sport := 0
	eport := 0
	portrage := common.GetSetting().Rmtserver.Localport

	if len(portrage) == 2 {
		sport = portrage[0]
		eport = portrage[1]
	} else {
		sport = 6001
		eport = 6200
	}
	return getFreePort(sport, eport)

}

func GetFreeProxyPort() int {
	sport := 0
	eport := 0
	portrage := common.GetSetting().Rmtserver.Proxyport

	if len(portrage) == 2 {
		sport = portrage[0]
		eport = portrage[1]
	} else {
		sport = 6201
		eport = 6400
	}
	return getFreePort(sport, eport)
}

func getFreePort(start, end int) int {
	for port := start; port < end; port++ {
		addr := fmt.Sprint("localhost:", port)
		if ln, err := net.Listen("tcp", addr); err == nil {
			ln.Close()
			return port
		}
	}
	return 0
}

func addRmtConnect(rmt RmtConnect) {
	lock.Lock()
	defer lock.Unlock()
	rmtclients[rmt.Ip] = &rmt
}

func setRmtConnectListener(port int, listener *net.TCPListener) {
	lock.Lock()
	defer lock.Unlock()
	for ip, rmt := range rmtclients {
		if rmt.Localport == port {
			rmtclients[ip].Locallistener = listener
			return
		}
	}
}

func delRmtConnect(ip string) {
	lock.Lock()
	defer lock.Unlock()
	for _, client := range rmtclients[ip].Clients {

		if client.Localconnert != nil {
			client.Localconnert.Close()
		}
		if client.Proxyconnect != nil {
			client.Proxyconnect.Close()
		}
		client.proxylistener.Close()
	}

	rmtclients[ip].Locallistener.Close()
	delete(rmtclients, ip)
}

func DelFreeRmtConnect() {

	freelist := make(map[int]string)
	for {
		for k := range freelist {
			delete(freelist, k)
		}

		time.Sleep(time.Minute * 1)
		lock.RLock()

		i := 0
		for ip := range rmtclients {
			if !findRtuClientFromIp(ip) {
				freelist[i] = ip
				i++
			}
		}
		lock.RUnlock()
		for _, ip := range freelist {
			delRmtConnect(ip)
		}
	}
}

func findRtuClientFromIp(ip string) bool {
	for _, rtu := range rtuclients {
		if ip == rtu.Ip {
			return true
		}
	}
	return false

}

func delRmtClient(localport, proxyport int) {
	lock.Lock()
	defer lock.Unlock()
	for ip, client := range rmtclients {
		if client.Localport == localport {
			var clts []RmtClient
			for _, c := range client.Clients {
				if c.Proxyport == proxyport {
					if c.proxylistener != nil {
						c.proxylistener.Close()
					}
				} else {
					clts = append(clts, c)
				}
			}
			rmtclients[ip].Clients = clts
			break
		}
	}
}

func findRmtClient(ip string) bool {
	lock.RLock()
	defer lock.RUnlock()
	for i := range rmtclients {
		if i == ip {
			return true
		}
	}
	return false
}

func findRmtClientFromLocalport(localport int) RmtConnect {
	lock.RLock()
	defer lock.RUnlock()
	for _, rmt := range rmtclients {
		if rmt.Localport == localport {
			return *rmt
		}
	}
	return RmtConnect{}
}

func findRmtLocalConnect(localport, proxyport int) net.Conn {
	lock.RLock()
	defer lock.RUnlock()
	for ip, rmt := range rmtclients {
		if rmt.Localport == localport {
			for i, rc := range rmtclients[ip].Clients {
				if rc.Proxyport == proxyport {
					return rmtclients[ip].Clients[i].Localconnert
				}
			}
		}
	}
	return nil
}

func addProxyRmtClient(localport int, rc RmtClient) {
	lock.Lock()
	defer lock.Unlock()
	for ip, rmt := range rmtclients {
		if rmt.Localport == localport {
			rmtclients[ip].Clients = append(rmtclients[ip].Clients, rc)

			return
		}
	}
}

func setProxyListen(localport int, proxyport int, listener *net.TCPListener) {
	lock.Lock()
	defer lock.Unlock()
	for ip, rmt := range rmtclients {
		if rmt.Localport == localport {
			for i, rc := range rmtclients[ip].Clients {
				if rc.Proxyport == proxyport {
					rmtclients[ip].Clients[i].proxylistener = listener
					return
				}
			}
		}
	}
}

func setProxyConect(localport int, proxyport int, conn net.Conn) {
	lock.Lock()
	defer lock.Unlock()
	for ip, rmt := range rmtclients {
		if rmt.Localport == localport {
			for i, rc := range rmtclients[ip].Clients {
				if rc.Proxyport == proxyport {
					rmtclients[ip].Clients[i].Proxyconnect = conn
					rmtclients[ip].Clients[i].Proxyip = conn.RemoteAddr().String()
					return
				}
			}
		}
	}
}
