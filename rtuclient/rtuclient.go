package main

import (
	"config"
	"fmt"
	"os"
	"os/signal"
	"rtuclient/common"
	"syscall"
	"time"
)

const VERSION = "1.0.0"

var recivebuff = make([]byte, 1024)

func main() {

	if len(os.Args) > 1 {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	config.SavePid(config.RtuclientPid)
	//defer func() {
	//	config.RemovePid(config.RtuclientPid)
	//}()

	//优雅的退出
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-c
		config.RemovePid(config.RtuclientPid)
		config.Led2Off()
		//config.SetLed(2, false)
		os.Exit(0)
	}()

	//go config.StartPprof(":6060")
	//fmt.Println(config.GetPid(config.RtuclientPid))

	common.InitSetting()
	//setting = config.GetSetting()

	//断线重连

	//go StartPprof()

	go NtpClock()

	go HeartbeatSignal()

	//主循环，处理请求
	for {

		if !isConnected() {
			//fmt.Println("con..")
			go config.Led2Off()
			Connect()
			continue
		} else {
			go config.Led2On()
		}

		//recivebuff := make([]byte, 1024)
		n, err := conn.Read(recivebuff)
		if err != nil {
			setConnected(false)

			conn.Close()
			continue
		}

		ProcRecive(recivebuff[:n])
		//fmt.Print(".")

		time.Sleep(100 * time.Millisecond)
	}

}
