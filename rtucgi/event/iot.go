package event

import (
	"config"
	"encoding/json"
	"main/common"
	"net/http"
	"strconv"
)

func Iots(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("iot")
	if r.Method == "GET" {
		striot := common.GetAnyQuery("iot")

		cfg, ok := config.LoadIotReadConfig()
		if !ok {
			common.ResponseFailInfo(w, "读取传感器信息失败！")
			return
		}

		if striot == "" {
			common.ResponseInfo(w, config.GetIots(cfg))
		} else {
			iot, errc := strconv.Atoi(striot)
			if errc != nil {
				common.ResponseFailInfo(w, "读取传感器信息失败！")
				return
			}
			if config.FindIot(cfg, iot) {
				common.ResponseInfo(w, config.GetIot(cfg, iot))
			} else {
				common.ResponseFailInfo(w, "读取传iot不存在！")
			}

		}
	} else if r.Method == "POST" {
		cfg, ok := config.LoadIotReadConfig()
		if !ok {
			common.ResponseFailInfo(w, "读取传感器信息失败！")
			return
		}

		iot, errc := strconv.Atoi(common.GetAnyQuery("iot"))
		if errc != nil {
			common.ResponseFailInfo(w, "参数id错误！")
			return
		}

		para := common.GetAnyQuery("para")
		switch para {
		case config.CMD_PARA_ADD:
			if !config.FindIot(cfg, iot) {
				iotcfg := config.IotConfig{}
				errd := json.NewDecoder(r.Body).Decode(&iotcfg)
				if errd != nil {
					common.ResponseFailInfo(w, "获取传感器信息失败！")
					return
				}

				cfg.Iots[iot] = iotcfg
				if config.SaveIotReadConfig(cfg) {
					common.ResponseSuccessNoData(w)
				} else {
					common.ResponseFailInfo(w, "添加失败！")
					return
				}

			} else {
				common.ResponseFailInfo(w, "id已经存在！")
				return
			}

		case config.CMD_PARA_DEL:
			if config.FindIot(cfg, iot) {
				delete(cfg.Iots, iot)
				if config.SaveIotReadConfig(cfg) {
					common.ResponseSuccessNoData(w)
				} else {
					common.ResponseFailInfo(w, "删除失败！")
					return
				}
			} else {
				common.ResponseFailInfo(w, "id不存在！")
				return
			}

		case config.CMD_PARA_EDIT:
			if config.FindIot(cfg, iot) {
				iotcfg := config.IotConfig{}
				errd := json.NewDecoder(r.Body).Decode(&iotcfg)
				if errd != nil {
					common.ResponseFailInfo(w, "获取传感器信息失败！")
					return
				}
				cfg.Iots[iot] = iotcfg
				if config.SaveIotReadConfig(cfg) {
					common.ResponseSuccessNoData(w)
					return
				} else {
					common.ResponseFailInfo(w, "修改失败！")
					return
				}
			} else {
				common.ResponseFailInfo(w, "id不存在！")
				return
			}

		}

	}
}

/*
func getIots() []config.IotConfig {

	cfg, ok := config.LoadIotReadConfig()
	if ok {
		return config.GetIots(cfg)

	}
	return nil

}

func getIot(id string) (config.IotConfig, bool) {

	//var result common.IotConfig
	for _, iot := range getIots() {
		if iot.Id == id {
			return iot, true
		}
	}
	return config.IotConfig{}, false

}
*/
