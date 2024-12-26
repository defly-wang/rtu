package event

import (
	"config"
	"encoding/json"
	"main/common"
	"net/http"
)

func Webapi(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		cfg, ok := config.LoadIotReadConfig()
		if ok {
			common.ResponseInfo(w, cfg.Webapi)
		} else {
			common.ResponseFailInfo(w, "读取配置失败！")
		}

	} else if r.Method == "POST" {
		webapi := config.WebapiConfig{}
		err := json.NewDecoder(r.Body).Decode(&webapi)
		if err != nil { //	result.Status = http.StatusInternalServerError
			common.ResponseFailInfo(w, "读取输入参数失败！")
			return
		}

		cfg, ok := config.LoadIotReadConfig()
		if ok {
			cfg.Webapi = webapi
			if config.SaveIotReadConfig(cfg) {
				common.ResponseSuccessNoData(w)
			} else {
				common.ResponseFailInfo(w, "保存配置失败！")
			}

		} else {
			common.ResponseFailInfo(w, "读取旧配置失败！")
		}
	}

}
