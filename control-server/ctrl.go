package main

import (
	"config"
	"control-server/common"
	"encoding/json"
	"net"
	"time"
)

func CtlSocket() {
	ctlserver, _ := net.ResolveTCPAddr("tcp4", common.GetCtlAddr())
	ctllisten, err := net.ListenTCP("tcp", ctlserver)
	if err != nil {
		return
	}
	defer ctllisten.Close()

	for {

		conn, err := ctllisten.Accept()
		if err != nil {
			//fmt.Println("accept:", err.Error())
			//删除该连接
			continue
		}
		clientid := saveCtlConnect(conn)
		go haneleCtl(clientid, conn)

		time.Sleep(10 * time.Millisecond)

	}
}

func haneleCtl(clientid uint, conn net.Conn) {
	defer func() {
		conn.Close()
		delCtlClient(clientid)
	}()

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			//fmt.Println("Read:", err.Error())
			//删除该连接
			return
		}

		CtlProcRecive(clientid, buf[:n])
	}
}

func CtlProcRecive(clientid uint, recive []byte) {
	var cmd config.Cmd

	jerr := json.Unmarshal(recive, &cmd)
	if jerr != nil {
		//fmt.Println(jerr.Error())
		//do noting
		return
	}
	switch cmd.Cmd {
	case config.CMD_RTUClients:
		sendcmd := config.Cmd{
			Cmd:  config.CMD_RTUClients,
			Dec:  getCtlClient(clientid).Ip,
			Data: GetRtuClients(),
		}
		SendDataJson(getCtlClient(clientid).Conn, sendcmd)
	default:
		//转发
		//增加发送源头和类型
		cmd.Src = getCtlClient(clientid).Ip
		cmd.Type = config.ClientTypeCTL
		decid := findRtuClient(cmd)
		if decid != 0 {
			//找到，转发到目标
			SendDataJson(getRtuClient(decid).Conn, cmd)
			//rtuclients[decid].Conn.Write(recive)
		} else {

			SendDataJson(getCtlClient(clientid).Conn, errorCmd(cmd.Src))
		}
	}
}

func saveCtlConnect(conn net.Conn) uint {
	fd, _ := conn.(*net.TCPConn).File()
	//fmt.Println(conn.RemoteAddr().String())
	var connect = CtlConnect{
		Conn: conn,
		Ip:   conn.RemoteAddr().String(),
	}
	intfd := fd.Fd()
	lock.Lock()
	defer lock.Unlock()

	ctlclients[uint(intfd)] = &connect
	return uint(intfd)
}

func findCtlClient(ip string) uint {
	lock.RLock()
	defer lock.RUnlock()

	for id, ctl := range ctlclients {
		if ip == ctl.Ip {
			return id
		}
	}
	return 0
}

func delCtlClient(id uint){
	lock.Lock()
	defer lock.Unlock()
	delete(ctlclients, id)
}


func getCtlClient(id uint) *CtlConnect{
	lock.RLock()
	defer lock.RUnlock()
	return ctlclients[id]
}