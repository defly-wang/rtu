package main

import (
	"config"
	"control-server/common"
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"github.com/gorilla/websocket"
)

type RtuConnect struct {
	Conn net.Conn `json:"-"`
	Ip   string   `json:"ip"`
	Mac  string   `json:"mac"`
	Unit string   `json:"unit"`
	Id   string   `json:"id"`
}

type CtlConnect struct {
	Conn net.Conn `json:"-"`
	Ip   string   `json:"ip"`
}

type WebConnect struct {
	Conn *websocket.Conn `json:"-"`
	Ip   string          `json:"ip"`
}

/*
	type ProxyConnect struct {
		Dtumac     string
		DirectPort uint
		DtuIp      string
		DtuConn    net.Conn
		//only 1
		DirectConn net.Conn
		DirectIp   string
	}
*/

type RmtClient struct {
	Proxyport     int              //代理端口
	proxylistener *net.TCPListener //代理监听
	Proxyconnect  net.Conn         //代理连接
	Localconnert  net.Conn         //接入客户 从rtuclient接入
	Proxyip       string
}

//var done = make(chan bool)

type RmtConnect struct {
	Ip            string           //rtu Ip(标识、dec)
	Remoteport    int              //rtu 服务端口
	Localport     int              //本地端口
	Locallistener *net.TCPListener //LocalServer接受连接时，生成代理端口并监听；发送给rtu，让rtu建立连接，将连接情况返回（Proxy connect Ip)
	Done          chan bool        `json:"-"` //标志
	Clients       []RmtClient      //连接的客户 每个客户一个连接，Local accept
}

var rtuclients = make(map[uint]*RtuConnect)
var ctlclients = make(map[uint]*CtlConnect)
var webclients = make(map[uint]*WebConnect)
var rmtclients = make(map[string]*RmtConnect)

//var proxyclients = make(map[uint]*config.ProxyConnect)

var lock sync.RWMutex

func main() {

	//rmtclients.delete
	common.InitSetting()

	go CtlWebSocket()

	go RtuSocket()

	go CtlSocket()

	go DelFreeRmtConnect()

	//go ProxyServer()

	select {}

}

func errorCmd(dec string) config.Cmd {
	return config.Cmd{
		Dec:  dec,
		Para: config.CMD_error, //"error",
		Data: "dec not find!",
	}
}

func SendDataJson(conn net.Conn, sendcmd config.Cmd) bool {
	if conn == nil {
		return false
	}

	bytes, err := json.Marshal(sendcmd)
	if err != nil {
		fmt.Println(err)
		return false
	}
	n, werr := conn.Write(bytes)
	if werr != nil {
		fmt.Println(err)
		return false
	}
	if n < len(bytes) {
		fmt.Println("send not complate!")
	}
	return true
}

func SendData(conn net.Conn, buff []byte) bool {
	if conn == nil {
		return false
	}

	n, err := conn.Write(buff)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if n < len(buff) {
		fmt.Println("send not complate!")
	}
	return true
}

func SendWebDataJson(conn *websocket.Conn, sendcmd config.Cmd) bool {
	if conn == nil {
		return false
	}

	err := conn.WriteJSON(sendcmd)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func SendWebData(conn *websocket.Conn, buff []byte) {
	if conn == nil {
		return
	}

	err := conn.WriteMessage(websocket.TextMessage, buff)
	if err != nil {
		fmt.Println(err)
	}
}

func GetRtuClients() []RtuConnect {
	var rtus []RtuConnect
	//var rtu Connect
	lock.RLock()
	defer lock.RUnlock()
	for _, rtu := range rtuclients {
		rtus = append(rtus, *rtu)
	}
	return rtus
}

func GetRtmClients() []RmtConnect {
	var rmts []RmtConnect
	lock.RLock()
	defer lock.RUnlock()
	for _, rmt := range rmtclients {
		rmts = append(rmts, *rmt)
	}
	return rmts
}

func GetRtmClientFromDec(dec string) *RmtConnect {
	lock.RLock()
	defer lock.RUnlock()
	for _, client := range rmtclients {
		if client.Ip == dec {
			return client
		}
	}
	return nil
}

/*
func GetDataFromCmd(buff []byte, data any) {

	//variableType := reflect.TypeOf(data)
	//variableType reflect.Type
	//data := reflect.New(variableType)
	dc := viper.New()
	dc.SetConfigType("json")
	//vv.ReadConfig()
	err := dc.ReadConfig(strings.NewReader(string(buff)))
	if err != nil {
		fmt.Println(err)
		return
	}
	erru := dc.Sub("data").Unmarshal(&data)
	if erru != nil {
		fmt.Println("Unmarshal:", erru)
	}
	//fmt.Println(data)
	//return data
}
*/
