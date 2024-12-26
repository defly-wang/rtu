package send

import (
	"config"
	"encoding/json"
	"fmt"
	"net"
)

var Conn net.Conn

func SocketInit() bool {

	//conn.RemoteAddr()
	conn, err := net.Dial("tcp", "192.168.1.128:5001")
	if err != nil {
		fmt.Println("net conn error")
		return false
	}

	Conn = conn
	go NetReadIng()
	//defer conn.Close()
	//Conn = nil
	//fmt.Println(conn)
	//conn.Write()
	//time.Tick(0)

	//conn.LocalAddr()
	return true
	// defer conn.Close()
}

func NetSend(buff []byte) bool {
	if _, err := Conn.Write(buff); err != nil {
		return false
	}
	return true
}

func NetReadIng() {
	buff := make([]byte, 1024)
	for {
		//time.Sleep(10 * time.Millisecond)

		len, err := Conn.Read(buff)
		if err != nil {
			//断网
			if err.Error() == "EOF" {
				fmt.Println("-----")
				Conn = nil
				break
			}
		}
		//收到数据
		if len > 0 {
			recivebuff := make([]byte, len)
			/*
				for i := 0; i < len; i++ {
					recivebuff[i] = buff[i]
				}
			*/
			copy(recivebuff, buff)

			fmt.Println("socket recive:", string(recivebuff))
		}

	}
}

func SendToAlarm(data config.IotInfo) bool {
	conn, err := net.Dial("unix", config.UNIX_ALARM_SOCKET)
	if err != nil {
		return false
	}
	defer conn.Close()

	jsondata, err := json.Marshal(data)
	if err != nil {
		return false
	}

	_, errr := conn.Write(jsondata)
	return errr == nil
}
