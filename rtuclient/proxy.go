package main

import (
	"fmt"
	"net"
	"rtuclient/common"
)

type DtuProxyConnect struct {
	TargetConn net.Conn
	SourceConn net.Conn
	SourceIp   string
	Connected  bool
	TargetPort uint
}

//var DtuProxyInfo DtuProxyConnect
/*
func GetSourceIp() string {
	if DtuProxyInfo.SourceConn != nil {
		return DtuProxyInfo.SourceConn.LocalAddr().String()
	}
	return ""
}
*/
func StartRemote(remoteport, proxyport int) (bool, string) {
	/*
		if DtuProxyInfo.Connected {
			return true
		}
	*/
	targatConn, ret := ConnectTargetServer(remoteport)
	if !ret {
		//DtuProxyInfo.Connected = false
		return false, ""
	}

	proxyConn, retp := ConnectProxy(proxyport)
	//fmt.Println("t ok")
	if !retp {
		//DtuProxyInfo.Connected = false
		return false, ""
	}
	//fmt.Println("p ok")

	//DtuProxyInfo.SourceIp = DtuProxyInfo.SourceConn.LocalAddr().String()
	//DtuProxyInfo.TargetPort = 22

	go copyData(targatConn, proxyConn)

	go copyData(proxyConn, targatConn)
	//fmt.Println("GO ok")

	//DtuProxyInfo.Connected = true

	return true, proxyConn.LocalAddr().String()

}

/*
func StopRemote() {

	DtuProxyInfo.TargetConn.Close()
	DtuProxyInfo.SourceConn.Close()
	DtuProxyInfo.TargetConn = nil
	DtuProxyInfo.SourceConn = nil
	DtuProxyInfo.Connected = false

	//go ConnectTargetcServer()
	//go ConnectProxy()

}
*/
func ConnectTargetServer(port int) (net.Conn, bool) {

	targetAddr := fmt.Sprint("localhost:", port)
	//var err error
	TargetConn, err := net.Dial("tcp", targetAddr)
	if err != nil {
		//DtuProxyInfo.TargetConn = nil
		//targetConnected = false
		return nil, false
	}
	return TargetConn, true

}

func ConnectProxy(port int) (net.Conn, bool) {
	//a :=
	proxyAddr := fmt.Sprint(common.GetServerAddress(), ":", port) //conn.RemoteAddr().String()
	//var err error
	SourceConn, err := net.Dial("tcp", proxyAddr)
	if err != nil {
		//DtuProxyInfo.SourceConn = nil
		return nil, false
	}
	return SourceConn, true
}

func copyData(source, dest net.Conn) {
	buf := make([]byte, 4096)
	for {
		if source == nil || dest == nil {
			return
		}
		n, err := source.Read(buf)
		if err != nil {
			//fmt.Printf("Error reading from source: %s\n", err)
			//continue
			break
		}
		//dest.SetWriteDeadline(time.Now().Add(10 * time.Second))
		_, err = dest.Write(buf[:n])
		if err != nil {
			//fmt.Printf("Error writing to dest: %s\n", err)
			break
			//continue
		}
	}
	defer func() {
		source.Close()
		dest.Close()
	}()

}
