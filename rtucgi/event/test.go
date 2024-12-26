package event

import (
	"config"
	"encoding/json"
	"main/common"
	"net/http"
)

func Test(w http.ResponseWriter, r *http.Request) {

	//fmt.Println("测试方式，显示在CGI服务中")
	//fmt.Println("mac=", common.GetMac())
	//w.Write([]byte("测试方式，显示在浏览器或测试工具中"))

	if r.Method == "GET" {
		common.ResponseInfo(w, "GET")

	} else if r.Method == "POST" {
		para := common.GetAnyQuery("para")

		switch para {
		case config.CMD_PARA_READ:
			iotcfg := config.IotConfig{}
			err := json.NewDecoder(r.Body).Decode(&iotcfg)
			if err != nil {
				common.ResponseFailInfo(w, "参数错误！")
				return
			}
			iotinfo, errr := config.TestReadSirial(iotcfg)
			if errr != nil {
				common.ResponseFailInfo(w, "读取测试错误！")
				return
			}
			common.ResponseInfo(w, iotinfo)

		case config.CMD_Mqtt:
			mqttinfo := config.MqttConfig{}
			err := json.NewDecoder(r.Body).Decode(&mqttinfo)
			if err != nil {
				common.ResponseFailInfo(w, "参数错误！")
				return
			}
			errm := config.TestMqtt(mqttinfo)
			if errm != nil {
				common.ResponseFailInfo(w, "MQTT服务器连接错误！"+errm.Error())
				return
			}
			common.ResponseSuccessNoData(w)
		case config.CMD_Webapi:
			webapi := config.WebapiConfig{}
			err := json.NewDecoder(r.Body).Decode(&webapi)
			if err != nil {
				common.ResponseFailInfo(w, "参数错误！")
				return
			}
			errm := config.TestWebapi(webapi)
			if errm != nil {
				common.ResponseFailInfo(w, "Webapi服务器连接错误！"+errm.Error())
				return
			}
			common.ResponseSuccessNoData(w)
		default:
			common.ResponseFailInfo(w, "参数错误！")
		}

	} else {
		common.ResponseFailInfo(w, "参数错误！")
		//user, _ := common.GetUserInfo()
		//common.ResponseInfo(w, user)
	}

}
