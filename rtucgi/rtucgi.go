package main

import (
	"config"
	"fmt"
	"main/common"
	"main/event"
	"net/http"
	"net/http/cgi"
	"os"
)

/*
	type ArgOptions struct {
		Help    bool `argv:"-h,--help"`
		Version bool `argv:"-v,--version"`
	}
*/
var Version = "1.0.1"

func main() {

	//var cmdop ArgOptions
	//argv.NewParser(&cmdop)
	//argv.Cmds()
	//argv.Argv()
	if len(os.Args) > 1 {
		fmt.Println(Version)
		os.Exit(0)
	}
	//argv
	//fmt.Println(os.Args)

	//fmt.Println("mac=", common.GetMac())
	cgi.Serve(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//fmt.Println("mac=", common.GetMac())
		//rut.cgi?action=login&token=aqwerqwerqwerasfdasdfasdf
		//login成功将返回token   data:"asdfasdfasqwrqw234"
		//其他模块可增加token参数

		//判断token是否有效 如不认证登陆状态，可去掉  login可不认证
		if common.GetQueryAction() != "login" && !common.IsTokenValid() {
			//或定义其他返回值 用以判断登陆过期或无效
			common.ResponseFailInfo(w, "认证无效")
			return
		}

		//rut.cgi?action=login
		//根据action执行不同的处理
		switch common.GetQueryAction() {
		case "login":
			event.Login(w, r)
		case "test":
			event.Test(w, r)
		case config.CMD_Iot:
			event.Iots(w, r)
		case config.CMD_Info:
			event.Sysinfo(w, r)
		case config.CMD_Status:
			event.Status(w, r)
		case config.CMD_Base:
			event.Base(w, r)
		case config.CMD_Mqtt:
			event.Mqtt(w, r)
		case config.CMD_Webapi:
			event.Webapi(w, r)
		case config.CMD_Power:
			event.Power(w, r)
		case config.CMD_Histiry:
			event.History(w, r)
		case config.CMD_Help:
			event.Help(w, r)

		default:
			common.ResponseFailInfo(w, "命令不支持！")
		}
	}))
}
