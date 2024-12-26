package main

import (
	"config"
	"control-server/common"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gorilla/websocket"
)

func CtlWebSocket() {
	http.HandleFunc("/", handleWeb)
	http.ListenAndServe(common.GetWebAddr(), nil)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWeb(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		//fmt.Println("close:::", ws)
		fmt.Println(err)
	}
	//conn := ws.NetConn()
	saveWenConnect(ws)

	//websocket.Conn = conn
	go handleMessage(ws)
}

func handleMessage(ws *websocket.Conn) {
	var cmd config.Cmd
	for {

		//recive := make([]byte, 1024)
		mType, recive, err := ws.ReadMessage()
		if err != nil {
			wsid := findWebClient(ws.RemoteAddr().String())
			if wsid != 0 {
				delWenConnect(wsid)
			}
			//主动断开

			//fmt.Println("aaaaaaaaaaaaaaaaa", wsid)
			//fmt.Println(err)
			return
		}
		//初始化
		cmd = config.Cmd{}
		if mType == websocket.TextMessage {
			errj := json.Unmarshal(recive, &cmd)
			if errj != nil {
				cmd.Para = config.CMD_paraerror
				cmd.Data = nil
				ws.WriteJSON(cmd)
				continue
			}
			//config.Cmd
		}

		/*
			sendcmd := config.Cmd{

			}
		*/

		sendcmd := config.Cmd{
			Cmd:  cmd.Cmd,
			Dec:  cmd.Src,
			Src:  cmd.Dec,
			Type: cmd.Type,
			Para: config.CMD_ok,
		}

		switch cmd.Cmd {
		case config.CMD_RTUClients:

			/*
				sendcmd := config.Cmd{
					Cmd:  config.CMD_RTUClients,
					Para: cmd.Para,
					//Dec:  ctlclients[clientid].Ip,
					Data: GetRtuClients(),
				}
			*/
			sendcmd.Data = GetRtuClients()
			ws.WriteJSON(sendcmd)

		case config.CMD_Remote:

			if cmd.Para == "get=all" {
				sendcmd.Data = GetRtmClients()
				ws.WriteJSON(sendcmd)
				continue
			}

			//rtu不错在ull}
			if findRtuClient(cmd) == 0 {
				sendcmd.Para = config.CMD_noexist
				ws.WriteJSON(sendcmd)
				continue
			}

			switch cmd.Para {
			case config.CMD_PARA_GET:
				sendcmd.Data = GetRtmClientFromDec(cmd.Dec)
				ws.WriteJSON(sendcmd)

			case config.CMD_PARA_ADD:
				//s := cmd.Data.(string)
				//存在
				if findRmtClient(cmd.Dec) {
					sendcmd.Para = config.CMD_existed
					//cmd.Data = nil
					ws.WriteJSON(sendcmd)
					continue
				}

				strport := cmd.Data
				//fmt.Println(strport)

				if reflect.TypeOf(strport) != reflect.TypeOf("200") {
					//fmt.Println("!str")
					sendcmd.Para = config.CMD_paraerror
					//cmd.Para = config.CMD_paraerror
					//cmd.Data = nil
					ws.WriteJSON(sendcmd)
					continue
				}

				remoteport, err := strconv.Atoi(strport.(string))
				if err != nil {
					//fmt.Println("ati")
					//cmd.Para = config.CMD_paraerror
					//cmd.Data = nil
					sendcmd.Para = config.CMD_paraerror
					ws.WriteJSON(sendcmd)
					continue
				}

				if remoteport == 0 {
					sendcmd.Para = config.CMD_paraerror
					//cmd.Data = nil
					ws.WriteJSON(sendcmd)
					continue
				}

				freelocalport := GetFreeLocalPort()

				rmt := RmtConnect{
					Ip:         cmd.Dec,
					Localport:  freelocalport,
					Remoteport: remoteport,
					//Done:       make(chan bool),
				}

				go LocalServer(freelocalport)

				addRmtConnect(rmt)

				//sendcmd.Para = config.CMD_ok
				sendcmd.Data = rmt
				ws.WriteJSON(sendcmd)

			case config.CMD_PARA_DEL:
				//不存在
				if !findRmtClient(cmd.Dec) {
					sendcmd.Para = config.CMD_noexist
					//cmd.Data = nil
					ws.WriteJSON(sendcmd)
					continue
				}
				delRmtConnect(cmd.Dec)
				//删除连接
				//cmd.Para = config.CMD_ok
				//cmd.Data = nil
				ws.WriteJSON(sendcmd)

			default:
				sendcmd.Para = config.CMD_paraerror
				//cmd.Data = nil
				ws.WriteJSON(sendcmd)
			}
		default:
			//fmt.Println("web socket cmd=", cmd)

			//增加发送源头和类型
			cmd.Src = ws.NetConn().RemoteAddr().String()
			cmd.Type = config.ClientTypeWEB

			decid := findRtuClient(cmd)
			if decid != 0 {
				//找到，转发到目标
				//fmt.Println("web socket send to rtu ok!!!")
				SendDataJson(getRtuClient(decid).Conn, cmd)
			} else {
				//对于RTU 评估是否需要返回错误信息
				errcmd := config.Cmd{
					Cmd:  cmd.Cmd,
					Dec:  cmd.Src,
					Para: "error",
					Data: "dec not find!",
				}

				ws.WriteJSON(errcmd)
				//SendData(ctlclients[clientid].Conn, errcmd)
			}
		}

	}
}

func saveWenConnect(conn *websocket.Conn) uint {
	fd, _ := conn.NetConn().(*net.TCPConn).File()
	//fmt.Println(conn.RemoteAddr().String())

	var connect = WebConnect{
		Conn: conn,
		Ip:   conn.RemoteAddr().String(),
	}

	intfd := fd.Fd()
	//fmt.Println("save web con,conn=", connect, "id=", intfd)
	lock.Lock()
	defer lock.Unlock()

	webclients[uint(intfd)] = &connect
	return uint(intfd)
}

func delWenConnect(id uint) {
	lock.Lock()
	defer lock.Unlock()

	delete(webclients, id)
}

/*
	func findCtlClientByws(ws *websocket.Conn) uint {
		lock.RLock()
		for id, web := range webclients {
			if ws.RemoteAddr().String() == web.Ip {
				lock.RUnlock()
				return id
			}
		}
		lock.RUnlock()
		return 0
	}
*/
func findWebClient(ip string) uint {
	lock.RLock()
	defer lock.RUnlock()
	for id, web := range webclients {
		if ip == web.Ip {

			return id
		}
	}
	return 0
}
