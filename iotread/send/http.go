package send

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"iotread/common"
	"net/http"
)

func HttpSend(jsondata []byte) bool {

	client := &http.Client{}
	bodyReader := bytes.NewReader(jsondata)
	//a:= common.CFG.Webapi.Url
	//"http://192.168.1.128:8080//webinsert"
	request, err := http.NewRequestWithContext(context.Background(), http.MethodPost, common.CFG.Webapi.Url, bodyReader)
	if err != nil {
		return false
	}
	request.Header.Set("Content-Type", "application/json")
	if common.CFG.Webapi.Token != "" {
		request.Header.Add("Authorization", fmt.Sprint("Bearer ", common.CFG.Webapi.Token))
	}

	//request.
	resp, err1 := client.Do(request)
	if err1 != nil {
		//fmt.Println(err)
		return false
	}
	_, err2 := io.ReadAll(resp.Body)
	if err2 != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
	//fmt.Println("statusCode,Body: ", resp.StatusCode, string(body))
	//return true
}
