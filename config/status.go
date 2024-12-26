package config

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
)

type Status struct {
	Starttime   MyTime `json:"starttime"`
	Count       uint64 `json:"count"`
	Readcount   uint64 `json:"readcount"`
	Mqttcount   uint64 `json:"mqttcount"`
	Webapicount uint64 `json:"webapicount"`
	Runing      bool   `json:"runing"`
	Version     string `json:"version"`
}

const (
	IotreadPid        = "/var/run/iotread.pid"
	RtuclientPid      = "/var/run/rtuclient.pid"
	UNIX_SOCKET_NAME  = "@IOTREAD_SOCKET"
	UNIX_ALARM_SOCKET = "@ALARM_SOCKET"
)

func SavePid(file string) {
	pid := os.Getpid()

	// 将PID转换为字符串
	pidStr := strconv.Itoa(pid)

	// 将PID保存到文件
	if err := os.WriteFile(file, []byte(pidStr), 0644); err != nil {
		fmt.Println("Please Run in Root:", err)
		os.Exit(1)
	}
}

func GetPid(file string) int {
	if buff, err := os.ReadFile(file); err == nil {
		if ret, errc := strconv.Atoi(string(buff)); errc == nil {
			return ret
		}
	}
	return 0
}

func RemovePid(file string) {
	if os.Remove(file) != nil {
		fmt.Println("Remove pid file err.")
	}
}

func GetRtuRunStatus() Status {
	status := Status{}
	conn, err := net.Dial("unix", UNIX_SOCKET_NAME)
	if err != nil {
		//fmt.Println(err.Error())
		return status
	}
	defer conn.Close()

	buf := make([]byte, 512)
	n, errr := conn.Read(buf)
	if errr != nil {
		return status
	}
	errj := json.Unmarshal(buf[:n], &status)
	if errj != nil {
		return status
	}

	return status
}
