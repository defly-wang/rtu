package main

import (
	"config"
	"encoding/json"
	"fmt"
	"net"
	"rtuclient/common"
	"sync"
	"time"
)

var (
	conn           net.Conn = nil
	connected      bool     = false
	connectedMutex sync.RWMutex
)

func Connect() {

	connectretrydelay := common.GetConnectRetrylDelay()
	address := fmt.Sprint(common.GetSetting().Server.Address, ":", common.GetSetting().Server.Port)
	tcpServer, _ := net.ResolveTCPAddr("tcp4", address)

	var err error
	conn, err = net.DialTCP("tcp", nil, tcpServer)

	if err != nil {

		//fmt.Println(err.Error())
		time.Sleep(time.Duration(connectretrydelay) * time.Second)
		//continue
	} else {
		//fmt.Println("conned=", conn)
		setConnected(true)

		//发送心跳包
		//_, err := conn.Write(MakeHS())
		if !SendSignal(conn) {
			setConnected(false)
			conn.Close()
		}
		//defer conn.Close()
	}

}

func SendData(sendcmd config.Cmd) {
	//cichu

	bytes, err := json.Marshal(sendcmd)
	if err != nil {
		return
	}

	_, errc := conn.Write(bytes)
	if errc != nil {
		fmt.Println(errc)
	}
}

func isConnected() bool {
	connectedMutex.RLock()
	isconnected := connected
	connectedMutex.RUnlock()
	return isconnected
}

func setConnected(isconnected bool) {
	connectedMutex.Lock()
	connected = isconnected
	connectedMutex.Unlock()
}
