package main

import (
	"config"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type RainDataType struct {
	DayRain  int32
	Rain24   int32
	HourRain int32
	Datas    [1440]int32
	//Begin    time.Time
	Now time.Time
}

type AlarmType struct {
	AlarmValue   int32
	ReAlarmDelay int32 //秒
}

var RainData RainDataType

// 数据初始化
func initData() {
	RainData.DayRain = 0
	RainData.Rain24 = 0
	RainData.HourRain = 0
	for i := 0; i < len(RainData.Datas); i++ {
		RainData.Datas[i] = 0
	}
	t := time.Now()
	//RainData.Begin = t.Add(-24 * time.Hour).Truncate(time.Minute)
	RainData.Now = t.Truncate(time.Minute)
}

func insertData(data config.IotInfo) {

	newtime := time.Time(data.Time).Truncate(time.Minute)

	//diff := newtime.Sub(RainData.Now)
	minutus := int(newtime.Sub(RainData.Now).Minutes())
	fmt.Println("diff=", minutus)

	//超过1天
	if minutus > 1440 {
		initData()
		RainData.Now = newtime
		return
	}
	//第一次执行有几率min=0
	if minutus == 0 {
		RainData.Now = newtime
		return
	}

	//移位
	for i := len(RainData.Datas) - 1; i >= minutus; i-- {
		RainData.Datas[i] = RainData.Datas[i-minutus]
	}

	//平均分配
	minutusraun := int32(int64(data.Foffset*10) / int64(minutus))
	for i := 0; i < minutus; i++ {
		RainData.Datas[i] = minutusraun
	}

	RainData.Now = newtime

	RainData.HourRain = countHourRain()
	RainData.Rain24 = countRain24()
	RainData.DayRain = countDayRain()
}

// 计算日雨量(0点)
func countDayRain() int32 {

	now := RainData.Now
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	//diff := now.Sub(midnight) //.Truncate(time.Hour))
	minutus := int(now.Sub(midnight).Minutes())
	fmt.Println("midnight min=", minutus)

	var sum int32 = 0
	for i := 0; i < minutus; i++ {
		sum += RainData.Datas[i]
	}
	return sum
}

// 计算24小时雨量
func countRain24() int32 {
	var sum int32 = 0
	for i := 0; i < len(RainData.Datas); i++ {
		sum += RainData.Datas[i]
	}
	return sum
}

// 计算小时雨量
func countHourRain() int32 {
	var sum int32 = 0
	for i := 0; i < 60; i++ {
		sum += RainData.Datas[i]
	}
	return sum
}

func main() {

	//rain24 := 0.0
	initData()

	lisunix, err := net.Listen("unix", config.UNIX_ALARM_SOCKET)
	if err != nil {
		//fmt.Println(err)
		return
	}
	defer lisunix.Close()

	buff := make([]byte, 512)

	for {
		conn, err := lisunix.Accept()
		if err != nil {
			time.Sleep(10 * time.Millisecond)
			//fmt.Println(err.Error())
			continue
		}
		//发送后结束
		conn.SetDeadline(time.Now().Add(time.Second))
		redlen, erread := conn.Read(buff)
		if erread != nil {
			conn.Close()
			continue
		}
		data := config.IotInfo{}

		errj := json.Unmarshal(buff[:redlen], &data)
		if errj != nil {
			conn.Close()
			continue

		}
		insertData(data)
		fmt.Println("data", data)
		fmt.Println("foffset", data.Foffset)
		//fmt.Println("value", data.Value)
		fmt.Println("RainData:", RainData)
		//fmt.Println(data)

		time.Sleep(300 * time.Millisecond)
		conn.Close()
	}

}
