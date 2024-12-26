package event

import (
	"config"
	"main/common"
	"net/http"
)

func Sysinfo(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		common.ResponseInfo(w, config.GetSysInfo())
	} else {
		common.ResponseErrorInfo(w)
	}

}
