package config

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

const (
	Network    = "unix"
	Servername = "@phone_server"
	SendHead   = "HTS"
	ReviceHead = "HTR"

	//版本号
	CMD_GET_ATI = 51
	//模块的IMEI码
	CMD_GET_IMEI = 52
	//获取SIM卡的IMSI码
	CMD_GET_IMSI = 101
	//获取信号强度
	CMD_GET_CSQ = 102
	//获取运营商信息
	CMD_GET_COPS = 103
	//获取SIM卡的ICCID标识
	CMD_GET_CCID = 104
	//获取基站的LAC和CI
	CMD_GET_CREG = 111
	//获取基站定位的位置信息，经纬度
	CMD_AT_QLBS = 452

	//获取服务程序版本
	CMD_GET_VERSION = 911

	//发送模块复位命令
	CMD_RESET = 901

	CMD_OK = 0
)

func Get4Ginfo(cmd int) string {
	conn, err := net.Dial(Network, Servername)
	if err != nil {
		//fmt.Println(err.Error())
		return ""
	}
	defer conn.Close()

	sendbuf := make([]byte, 7)
	copy(sendbuf, SendHead)

	sendbuf[3] = byte(cmd >> 8)
	sendbuf[4] = byte(cmd)

	conn.Write(sendbuf)
	//time.Sleep(100 * time.Millisecond)
	conn.SetReadDeadline(time.Now().Add(300 * time.Millisecond))

	reader := bufio.NewReader(conn)
	recivebuff := make([]byte, 128)
	//for {

	n, err := reader.Read(recivebuff)
	if err != nil {
		return ""
	}
	//break

	if n < 7 { //太短
		return ""
	}
	if string(recivebuff[:3]) != ReviceHead { //头错误
		return ""
	}

	if int(recivebuff[3])<<8+int(recivebuff[4]) != cmd { //CMD字错误
		return ""
	}
	if int(recivebuff[5])<<8+int(recivebuff[6]) != CMD_OK { //返回错误
		return ""
	}

	//转换,截取 0x00
	ret := string(recivebuff[7 : 7+strlen(recivebuff[7:])])
	//fmt.Println("ret=", ret, "len=", len(ret))
	return ret

}

func strlen(buf []byte) int {

	for i, v := range buf {
		if v == 0 {
			return i
		}
	}
	return len(buf)
}

func Maintest() {

	fmt.Println("ATI", Get4Ginfo(CMD_GET_ATI))
	time.Sleep(200 * time.Millisecond)
	fmt.Println("CCID", Get4Ginfo(CMD_GET_CCID))
	time.Sleep(200 * time.Millisecond)
	fmt.Println("COPS", Get4Ginfo(CMD_GET_COPS))
	time.Sleep(200 * time.Millisecond)
	fmt.Println("CSQ", Get4Ginfo(CMD_GET_CSQ))
	time.Sleep(200 * time.Millisecond)
	fmt.Println("IMEI", Get4Ginfo(CMD_GET_IMEI))
	time.Sleep(200 * time.Millisecond)
	fmt.Println("IMSI", Get4Ginfo(CMD_GET_IMSI))
	time.Sleep(200 * time.Millisecond)
	fmt.Println("CREG", Get4Ginfo(CMD_GET_CREG))
	time.Sleep(200 * time.Millisecond)
	fmt.Println("QLBS", Get4Ginfo(CMD_AT_QLBS))

	time.Sleep(200 * time.Millisecond)
	fmt.Println("Version=", Get4Ginfo(CMD_GET_VERSION))
}
