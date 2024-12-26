package main

import (
	"config"
	"control-server/common"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

func RtuSocket() {
	rtuserver, _ := net.ResolveTCPAddr("tcp4", common.GetRtuAddr())
	rtulisten, err := net.ListenTCP("tcp", rtuserver)
	if err != nil {
		return
	}
	defer rtulisten.Close()

	for {
		conn, err := rtulisten.Accept()
		if err != nil {
			fmt.Println("accept:", err.Error())
			//删除该连接
			continue
		}
		clientid := saveRtuConnect(conn)
		go haneleRtu(clientid, conn)

		time.Sleep(100 * time.Millisecond)

	}
}

func haneleRtu(clientid uint, conn net.Conn) {
	defer func() {
		conn.Close()

		delRtuClient(clientid)
	}()

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			//fmt.Println("Read:", err)
			//删除该连接
			return
		}

		RtuProcRecive(clientid, buf[:n])

	}
}

func RtuProcRecive(clientid uint, recive []byte) {

	var cmd config.Cmd

	jerr := json.Unmarshal(recive, &cmd)
	if jerr != nil {
		//log.Printf("err=%v", jerr.Error())
		return
	}
	//fmt.Println(cmd)

	switch cmd.Cmd {

	case config.CMD_HeartbeatSignal:

		//存MAC TYPE等信息
		var hd = config.HeartbeatData{}

		errhd := config.GetDataFromCmd(recive, &hd)
		if errhd == nil { //无错执行
			saveRtuConnectInfo(clientid, &hd)
		}
		//保存数据库，可用webapi

		//空，废弃
		rtuclt := getRtuClient(clientid)
		if rtuclt != nil {

			sendcmd := config.Cmd{
				Cmd:  config.CMD_HeartbeatSignal,
				Dec:  rtuclt.Ip,
				Data: config.CMD_ok,
			}

			SendDataJson(rtuclt.Conn, sendcmd)
		}

	case config.CMD_Remote:
		//从rtu返回连接情况 连接失败，删除连接
		//fmt.Println(cmd)
		if cmd.Para != config.CMD_ok {
			var data config.Proxy
			config.GetDataFromCmd(recive, &data)
			delRmtClient(data.Localport, data.Proxyport)
		}

	default:
		//转发

		switch cmd.Type {
		case config.ClientTypeCTL:
			ctlid := findCtlClient(cmd.Dec)
			if ctlid != 0 {
				//找到，转发到目标
				SendData(getCtlClient(ctlid).Conn, recive)

			} else {
				//对于RTU 评估是否需要返回错误信息
				rtuclt := getRtuClient(clientid)
				if rtuclt != nil {
					SendDataJson(getRtuClient(clientid).Conn, errorCmd(cmd.Src))
				}
				//clients[clientid].Conn.Write(errresult)
				//可返回notfaound
			}
		case config.ClientTypeWEB:
			//fmt.Println("rtu recive：cmd=", cmd)
			webid := findWebClient(cmd.Dec)

			if webid != 0 {
				//fmt.Println(recive)
				SendWebData(webclients[webid].Conn, recive)
				//webclients[webid].Conn.WriteJSON(cmd)
			}
		}

	}

}

func saveRtuConnect(conn net.Conn) uint {
	fd, _ := conn.(*net.TCPConn).File()
	//fmt.Println(conn.RemoteAddr().String())

	var connect = RtuConnect{
		Conn: conn,
		Ip:   conn.RemoteAddr().String(),
	}
	intfd := fd.Fd()
	lock.Lock()
	defer lock.Unlock()
	rtuclients[uint(intfd)] = &connect
	return uint(intfd)
}

func saveRtuConnectInfo(id uint, hd *config.HeartbeatData) {
	lock.Lock()
	defer lock.Unlock()
	//fmt.Println(hd.Mac)
	//log.Printf("id=%d", id)
	//log.Printf("hd=%+v", hd)
	//log.Printf("rtu=%+v", rtuclients[id])
	if hd == nil {
		return
	}
	if rtuclients[id] == nil {
		//需要重新写入
		//delete(rtuclients, id)
		//fmt.Println("nil")
		return
	}

	/*
		rtucon := rtuclients[id]
		if rtucon != nil {
			rtucon.Mac = hd.Mac
			rtucon.Id = hd.Id
			rtucon.Unit = hd.Unit
			rtuclients[id] = rtucon
		}

		fmt.Println(id)
	*/

	rtuclients[id].Mac = hd.Mac
	rtuclients[id].Id = hd.Id
	rtuclients[id].Unit = hd.Unit

}

func findRtuClient(cmd config.Cmd) uint {
	lock.RLock()
	defer lock.RUnlock()
	for id, rtu := range rtuclients {
		if cmd.Dec == rtu.Ip {
			return id
		}
	}
	return 0
}

func delRtuClient(id uint) {
	lock.Lock()
	defer lock.Unlock()
	//add
	if rtuclients[id] != nil {
		if rtuclients[id].Conn != nil {
			rtuclients[id].Conn.Close()
		}
	}
	delete(rtuclients, id)
}

func getRtuClient(id uint) *RtuConnect {
	lock.RLock()
	defer lock.RUnlock()
	return rtuclients[id]
}

/*
func GetHbDataFromCmd(buff []byte) config.HeartbeatData {

	data := config.HeartbeatData{}
	dc := viper.New()
	dc.SetConfigType("json")
	//vv.ReadConfig()
	err := dc.ReadConfig(strings.NewReader(string(buff)))
	if err != nil {
		fmt.Println(err)
		return data
	}
	erru := dc.Sub("data").Unmarshal(&data)
	if erru != nil {
		fmt.Println(erru)
	}
	return data
}
*/
