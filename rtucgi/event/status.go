package event

import (
	"config"
	"main/common"
	"net/http"
)

func Status(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		common.ResponseInfo(w, config.GetRtuRunStatus())
	} else {
		common.ResponseErrorInfo(w)
	}

}
