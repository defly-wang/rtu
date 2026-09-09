package event

import (
	"config"
	"encoding/json"
	"main/common"
	"net/http"
)

func Network(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		cfg, ok := config.GetNetworkConfig()
		if ok != nil {
			common.ResponseInfo(w, cfg)

		} else {
			common.ResponseFailInfo(w, "读取配置失败！")
		}

	} else if r.Method == "POST" {
		cfg := config.NetworkConfig{}
		err := json.NewDecoder(r.Body).Decode(&cfg)
		if err != nil { //	result.Status = http.StatusInternalServerError
			common.ResponseFailInfo(w, "读取输入参数失败！")
			return
		}

		if config.SaveNetworkConfig(cfg) {
			common.ResponseSuccessNoData(w)
		} else {
			common.ResponseFailInfo(w, "保存配置失败！")
		}

	}
}
