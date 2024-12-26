package common

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

type Response struct {
	Status uint   `json:"code"`
	Msg    string `json:"msg"`
	Data   any    `json:"data"`
}

// 分析Url aaa=???&bbb=??
func GetQuery() map[string]string {
	pairs := strings.Split(os.Getenv("QUERY_STRING"), "&")
	params := make(map[string]string)
	for _, pair := range pairs {
		split := strings.SplitN(pair, "=", 2)
		if len(split) != 2 {
			continue // 跳过不符合规则的键值对
		}
		params[split[0]] = split[1]
	}
	return params
}

// 一个
func GetAnyQuery(parameter string) string {
	return GetQuery()[parameter]
}

// 分析Url ?action=???
func GetQueryAction() string {
	return GetQuery()["action"]
}

// 严重错误输出
func ResponseErrorInfo(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("内部错误：访问出错！"))
}

// 成功输出
func ResponseInfo(w http.ResponseWriter, value any) {
	var respo Response
	respo.Status = http.StatusOK
	respo.Msg = "Success"
	respo.Data = value
	jsondata, err := json.Marshal(respo)
	if err != nil {
		//fmt.Println("ASD")
		ResponseErrorInfo(w)
	}
	w.Write(jsondata)
}

// 成功输出
func ResponseSuccessNoData(w http.ResponseWriter) {
	var respo Response
	respo.Status = http.StatusOK
	respo.Msg = "Success"
	respo.Data = nil
	jsondata, err := json.Marshal(respo)
	if err != nil {
		//fmt.Println("ASD")
		ResponseErrorInfo(w)
	}
	w.Write(jsondata)
}

// 成功输出
func ResponseFailInfo(w http.ResponseWriter, msg string) {
	var respo Response
	respo.Status = http.StatusBadRequest
	respo.Msg = msg
	jsondata, err := json.Marshal(respo)
	if err != nil {
		//fmt.Println("ASD")
		ResponseErrorInfo(w)
	}
	w.Write(jsondata)
}
