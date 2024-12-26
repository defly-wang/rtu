package event

import (
	"config"
	"encoding/json"
	"main/common"
	"net/http"
)

func Base(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		cfg, ok := config.LoadIotReadConfig()
		if ok {
			common.ResponseInfo(w, cfg.Common)
		} else {
			common.ResponseFailInfo(w, "读取配置失败！")
		}

	} else if r.Method == "POST" {
		base := config.CommonConfig{}
		err := json.NewDecoder(r.Body).Decode(&base)
		if err != nil { //	result.Status = http.StatusInternalServerError
			common.ResponseFailInfo(w, "读取输入参数失败！")
			return
		}

		cfg, ok := config.LoadIotReadConfig()
		if ok {
			cfg.Common = base
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
