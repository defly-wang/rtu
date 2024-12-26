package event

import (
	"config"
	"main/common"
	"net/http"
	"strings"
)

func Help(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		para := common.GetAnyQuery("para")
		switch para {
		case config.CMD_Mqtt:

			mqtthelp, ok := config.LoadHelpMqtt()
			if !ok {
				common.ResponseFailInfo(w, "读取mqtt缺省值失败！")
			}
			common.ResponseInfo(w, mqtthelp)
		case config.CMD_Webapi:
			webapi, ok := config.LoadHelpWebapi()
			if !ok {
				common.ResponseFailInfo(w, "读取mqtt缺省值失败！")
			}
			common.ResponseInfo(w, webapi)
		case "org":
			common.ResponseInfo(w, config.LoadHelpIotOrg())
		case "type":
			common.ResponseInfo(w, config.LoadHelpIotType())
		case config.CMD_PARA_READ:
			file := common.GetAnyQuery("file")
			if len(strings.TrimSpace(file)) != 0 {
				//if file == nil
				cfg, err := config.GetIotHelpConfig(file)
				if err != nil {
					common.ResponseFailInfo(w, "读取"+file+"缺省值失败！")
				}
				common.ResponseInfo(w, cfg)
			} else {
				common.ResponseFailInfo(w, "file参数错误！")
			}

		default:
			common.ResponseFailInfo(w, "参数错误！")

		}

	}
}
