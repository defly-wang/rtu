package event

import (
	"config"
	"main/common"
	"net/http"
	"strconv"
)

func History(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//para=day,one
		//iot=1,2
		//time=2024-1-1 0:0:0
		para := common.GetAnyQuery("para")

		time, err := config.StrToTime(common.GetAnyQuery("time"))
		//	t, err := time.Parse(format, para[1])
		if err != nil {
			common.ResponseFailInfo(w, "time参数错误！")
			return
		}

		iot, erriot := strconv.Atoi(common.GetAnyQuery("iot"))
		if erriot != nil {
			common.ResponseFailInfo(w, "iot参数错误！")
		}
		cfg, ok := config.LoadIotReadConfig()
		if !ok {
			common.ResponseFailInfo(w, "读取参数失败！")
		}

		switch para {
		case config.CMD_PARA_HIS_DAY:
			res, _ := config.LoadHistoryData(cfg, time, iot)
			common.ResponseInfo(w, res)
		case config.CMD_PARA_HIS_ONE:
			res, _ := config.LoadHistoryDataOne(cfg, time, iot)
			common.ResponseInfo(w, res)
		default:
			common.ResponseFailInfo(w, "para参数错误！")
		}

	}
}
