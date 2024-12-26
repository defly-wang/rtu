package event

import (
	"config"
	"encoding/json"
	"main/common"
	"net/http"
)

func Power(w http.ResponseWriter, r *http.Request) {
	powerset, ok := config.GetPowerSet()
	if !ok {
		common.ResponseErrorInfo(w)
		return
	}
	if r.Method == "GET" {

		para := common.GetAnyQuery("para")
		switch para {
		case config.CMD_PARA_INFO:
			powerinfo, err := config.ReadPowerInfo(powerset)
			if err != nil {
				common.ResponseFailInfo(w, "读取电源信息失败！")
				return
			}
			common.ResponseInfo(w, powerinfo)

		case config.CMD_PARA_DATA:
			powerdata, err := config.ReadPowerRunData(powerset)
			if err != nil {
				common.ResponseFailInfo(w, "读取电源运行数据参数失败！")
				return
			}
			common.ResponseInfo(w, powerdata)

		case config.CMD_PARA_SET:
			set, err := config.ReadPowerSet(powerset)
			if err != nil {
				common.ResponseFailInfo(w, "读取电源设置参数失败！")
				return
			}
			common.ResponseInfo(w, set)

		//case config.CMD_PARA_EDIT:
		default:

		}

	} else if r.Method == "POST" {
		ptc := config.PowerTypeCapacity{}
		err := json.NewDecoder(r.Body).Decode(&ptc)
		if err != nil { //	result.Status = http.StatusInternalServerError
			common.ResponseFailInfo(w, "读取输入参数失败！")
			return
		}
		if config.SetPowerBatterType(powerset, ptc.BatterType) != nil {
			common.ResponseFailInfo(w, "设置电池类型失败！")
			return
		}
		if config.SetPowerBatterCapacity(powerset, ptc.BatterCapacity) != nil {
			common.ResponseFailInfo(w, "设置电池容量失败！")
			return
		}
		common.ResponseSuccessNoData(w)

	}

}
