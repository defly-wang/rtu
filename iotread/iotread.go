package main

import (
	"config"
	"fmt"
	"iotread/common"
	"iotread/readiot"
	"iotread/save"
	"iotread/send"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
)

// #include <stdio.h>
// void sayhello(){
// printf("hello\n");
// }
// import "C"

//var UNIX_SOCKET_NAME = "@IOTREAD_SOCKET"

func main() {

	config.SavePid(config.IotreadPid)

	defer func() {
		config.RemovePid(config.IotreadPid)
	}()

	//优雅的退出
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-c
		config.RemovePid(config.IotreadPid)
		//return
		os.Exit(0)
	}()

	//str := time.Now().Format("2006-01-02 15:04:05")
	//fmt.Println(str)
	//读取config.toml文件
	//go config.StartPprof(":7070")

	//CFG,_ = config.LoadConfig()

	if !common.InitConfig() {
		fmt.Println("Read config file error")
		return
		//os.Exit(0)
	}

	//fmt.Println(config.CFG)
	//fmt.Println(config.GetDataDir())

	//初始化并连接mqtt,好像可以断线重连
	//可改称循环连接
	if common.CFG.Mqtt.Used {
		if !send.MqttInit() {
			fmt.Println("Init Mqtt error")
			//return
		}
		defer send.MqttClose()
	}
	//defer send.mtqqClient.Disconnect(300)

	//send.SocketInit()

	//定时任务
	thiscron := cron.New(cron.WithSeconds())
	//iot设备数量
	//fmt.Println(common.CFG.Common.Iotnumber)
	//此处可改为CFG.Iots循环
	for i := 1; i <= common.CFG.Common.Iotnumber; i++ {
		//fmt.Println(i)
		cronjob := common.FuncIntJob(i, CronTask)

		//各个iot读取的定时字符串
		_, err := thiscron.AddJob(common.CFG.Iots[i].Cron, cronjob)

		if err != nil {
			fmt.Println("now Cron error")
		}
		//taskmap[i] = entrtyid
		//fmt.Println("执行任务", i, "启动,任务id:", i, "\t启动时间:", time.Now())

	}

	//thiscron.Schedule()

	thiscron.Start()
	defer thiscron.Stop()

	//设置运行状态
	common.SetStatusRuning(true)

	//fmt.Println("--工作正常")
	//time.Sleep(1 * time.Minute)
	//启动 unix socket
	lisunix, err := net.Listen("unix", config.UNIX_SOCKET_NAME)
	if err != nil {
		//fmt.Println(err)
		return
	}
	defer lisunix.Close()

	for {
		conn, err := lisunix.Accept()
		if err != nil {
			time.Sleep(10 * time.Millisecond)
			//fmt.Println(err.Error())
			continue
		}
		//发送后结束
		conn.SetDeadline(time.Now().Add(time.Second))
		conn.Write(common.GetStatusJson())
		time.Sleep(300 * time.Millisecond)
		conn.Close()
		//fmt.Println("A client connected : " + conn.RemoteAddr().String())
		//go pipe(conn)
	}

	//可作为判断网络运行情况判断
	//time.Sleep(100 * time.Millisecond)
	//fmt.Println(time.Now())
	/*
		//能够自动连接
		if !mtqqClient.IsConnected() {
			InitMqtt()
		}
	*/

}

// 定时读取处理任务 id=Iot编号,任务号
func CronTask(iot int) bool {

	//保存读取次数
	common.StatusCount()
	//fmt.Println("begin")

	//config.OffLed()
	go config.LedBlink(600)

	//读3次
	data, err := readiot.ReadIotJson(iot, 3)
	if err != nil {
		go config.LedBlink(3000)
		//fmt.Println(err)
		return false
	}

	//保存读取成功次数
	common.StatusReadSuccess()

	//config.OffLed()

	save.Save(data)

	//err ==nil 正常，err ！=nil 返回code=203错误，无数据
	//{code=200 msg=ok data={id=,type=,value=,offset=,rawdate=,time=}}
	//可调整为发送数据库软件
	//jsondata := config.IotDataToJson(data)
	jsondata := config.IotDataToJsonResult(data, err)
	if jsondata == nil {
		return false
	}
	//fmt.Println("save ok")

	send.SendToAlarm(data)
	//fmt.Println("ADS")
	//do nothing
	//common.

	if common.CFG.Mqtt.Used {
		if send.MqttPublicInfo(common.CFG.Mqtt.Pubtoptic, string(jsondata)) {
			//保存mqtt发送成功次数
			common.StatusMqttSuccess()
		}

	}

	if common.CFG.Webapi.Used {
		if send.HttpSend(jsondata) {
			//保存webapi发送成功次数
			common.StatusWebapiSuccess()
		}
	}
	return true
}
