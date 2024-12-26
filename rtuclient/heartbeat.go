package main

import (
	"config"
	"encoding/json"
	"fmt"
	"net"
	"rtuclient/common"
	"time"
)

// var recivebuff = make([]byte, 1024)
var hstime = time.Now()

// 时间同步进程
func NtpClock() {
	for {
		config.NtpDate()
		time.Sleep(12 * time.Hour)
	}
}

func HeartbeatSignal() {
	//var hssignaldelay int
	hssignaldelay := common.GetHSSignalDelay()

	for {
		//fmt.Println(time.Now())
		//fmt.Println(config.GetSetting())
		if isConnected() {
			//fmt.Println("send hs")
			//_, err := conn.Write(MakeHS())
			if !SendSignal(conn) {
				//fmt.Println("hr err")
				setConnected(false)
				conn.Close()
			}
			if hstime.Add(time.Duration(hssignaldelay) * time.Minute * 2).Before(time.Now()) {
				fmt.Println(time.Now())
				setConnected(false)
				conn.Close()
			}
		}

		//超过两次未收到hs回应

		time.Sleep(time.Duration(hssignaldelay) * time.Minute)
	}

}

func SendSignal(conn net.Conn) bool {
	buff, err := json.Marshal(config.Cmd{
		Cmd:  config.CMD_HeartbeatSignal,
		Src:  conn.LocalAddr().String(),
		Type: config.ClientTypeRTU,
		Data: config.HeartbeatData{
			Id:   common.GetSetting().Id,
			Mac:  common.GetSetting().Mac,
			Unit: common.GetSetting().Unit,
			//Proxyclientip: SourceConn.LocalAddr().String(),
			//Port:          22,
		},
	})
	if err != nil {
		return false
	}
	n, errw := conn.Write(buff)
	//fmt.Println(n)
	return errw == nil && n == len(buff)
}
