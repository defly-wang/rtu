package event

import (
	"encoding/json"
	"main/common"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {

	//提交方法判断
	if r.Method != "POST" {
		common.ResponseErrorInfo(w)
		return
	}

	loginuser := common.User{}
	err := json.NewDecoder(r.Body).Decode(&loginuser)
	if err != nil { //	result.Status = http.StatusInternalServerError
		common.ResponseFailInfo(w, "获取登陆信息失败")
		return
	}

	sysuser, ok := common.GetUserInfo()
	if !ok {
		common.ResponseFailInfo(w, "读取登陆信息失败")
	}
	//fmt.Println(loginuser.Password)
	//fmt.Println(loginuser.Username)
	//fmt.Println(sysuser.Password)
	//fmt.Println(sysuser.Username)
	if sysuser.Username == loginuser.Username && sysuser.Password == loginuser.Password {
		//验证成功，生成token
		token := common.GenerateToken()
		common.ResponseInfo(w, token)
	} else {
		common.ResponseFailInfo(w, "登陆失败")
	}

}

func SavePwd(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		common.ResponseErrorInfo(w)
		return
	}

	loginuser := common.User{}
	err := json.NewDecoder(r.Body).Decode(&loginuser)
	if err != nil { //	result.Status = http.StatusInternalServerError
		common.ResponseFailInfo(w, "获取登陆信息失败")
		return
	}

	if common.SavePassword(loginuser.Password) {
		common.ResponseSuccessNoData(w)

	} else {
		common.ResponseFailInfo(w, "修改密码失败")
	}
}
