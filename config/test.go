package config

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jacobsa/go-serial/serial"
)

//import "config"

func TestReadSirial(iotcfg IotConfig) (IotInfo, error) {

	iotinfo := IotInfo{
		Id:    iotcfg.Id,   // GetIotid(iot),
		Type:  iotcfg.Type, //tGetIotType(iot),
		Time:  MyTime(time.Now()),
		Ratio: iotcfg.Ratio, //GetIotRatio(iot),
	}

	options := SetSerialOption(SerialSet{
		Com:        "/dev/ttyS" + iotcfg.Com,
		Baudrate:   iotcfg.Baudrate,
		Databits:   iotcfg.Databits,
		Stopbits:   iotcfg.Stopbits,
		Paritymode: iotcfg.Paritymode,
	})
	port, err := serial.Open(options)
	if err != nil {
		return iotinfo, errors.New("串口打开失败！")
	}
	defer port.Close()

	sendbuff, errh := hex.DecodeString(iotcfg.Read_buff)
	if errh != nil {
		return iotinfo, errors.New("读取指令格式无效！")
	}

	//if checkCRC16()
	if !CheckCRC(sendbuff, iotcfg.Crc) {
		return iotinfo, errors.New(fmt.Sprint("读取指令", iotcfg.Crc, "校验错误！"))
	}

	sendlen, errw := port.Write(sendbuff)

	if errw != nil {
		return iotinfo, errors.New("发送读取指令错误！")
	}
	if sendlen != iotcfg.Read_len {
		return iotinfo, errors.New("发送读取指令错误：发送长度不匹配！")
	}

	time.Sleep(200 * time.Millisecond)

	recivebuff := make([]byte, 256)
	reclen, errr := port.Read(recivebuff)
	if errr != nil {
		return iotinfo, errors.New("读取数据错误！")
	}
	if reclen < 3 {
		return iotinfo, errors.New(fmt.Sprint("读取数据错误:返回", reclen, "个字节！"))
	}

	if !CheckCRC(recivebuff[:reclen], iotcfg.Crc) {
		return iotinfo, errors.New(fmt.Sprint("读取数据", iotcfg.Crc, "校验错误,数据：", hex.EncodeToString(recivebuff)))
	}
	/*
		if iotcfg.Crc == CRC_CRC16 && !checkCRC16(recivebuff[:reclen]) {
			return iotinfo, errors.New("crc16 check Error")
		}
		if iotcfg.Crc == CRC_XOR && !checkXOR(recivebuff[:reclen]) {
			return iotinfo, errors.New("crcxor check Error")
		}
	*/

	iotinfo.Rawdata = hex.EncodeToString(recivebuff[:reclen])

	//var result uint64
	//计算需要重写（各种类型）
	compute(iotcfg, recivebuff, &iotinfo)

	return iotinfo, nil

}

func TestMqtt(mqttinfo MqttConfig) error {
	opts := mqtt.NewClientOptions()
	server := fmt.Sprint("tcp://", mqttinfo.Host, ":", mqttinfo.Port)
	opts.AddBroker(server)
	opts.SetClientID(randomString(8))
	opts.SetUsername(mqttinfo.Username) //config.GetMqttUser())
	opts.SetPassword(mqttinfo.Password)
	//超时时间，根据情况调整
	opts.SetConnectTimeout(300 * time.Millisecond)

	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return errors.New(fmt.Sprint("连接mqtt服务器失败：", token.Error()))
	}
	defer mqttClient.Disconnect(200)

	if token := mqttClient.Publish(mqttinfo.Pubtoptic, 0, true, randomString(16)); token.Wait() && token.Error() != nil {
		//fmt.Println(token.Error())
		return errors.New(fmt.Sprint("发布主题失败：", token.Error()))
	}

	return nil

}

func TestWebapi(api WebapiConfig) error {

	//data := IotInfo{}
	res := Result{
		Code: 200,
		Msg:  "OK",
		Data: IotInfo{},
	}
	jsondata, errj := json.Marshal(res)
	if errj != nil {
		return errors.New("内部错误：数据解析错误！")
	}

	bodyReader := bytes.NewReader(jsondata)
	//a:= common.CFG.Webapi.Url
	//"http://192.168.1.128:8080//webinsert"
	request, err := http.NewRequestWithContext(context.Background(), http.MethodPost, api.Url, bodyReader)
	if err != nil {
		return errors.New(fmt.Sprint("地址连接错误：", err.Error()))
	}
	request.Header.Set("Content-Type", "application/json")
	if api.Token != "" {
		request.Header.Add("Authorization", fmt.Sprint("Bearer ", api.Token))
	}

	//request.
	client := &http.Client{}
	resp, err1 := client.Do(request)
	if err1 != nil {
		return errors.New(fmt.Sprint("数据请求错误：", err1.Error()))
		//fmt.Println(err)
		//return err1
	}
	defer resp.Body.Close()
	//client.Do()
	_, err2 := io.ReadAll(resp.Body)
	if err2 != nil {
		return errors.New(fmt.Sprint("数据返回错误：", err2.Error()))
		//return err2
	}

	if resp.StatusCode != 200 {
		return errors.New("返回状态错误：提交不成功")
	}
	return nil
}

// 生成随即字符串
func randomString(length int) string {
	k := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, k)
	if err != nil {
		panic(err.Error())
	}
	return base64.StdEncoding.EncodeToString(k)
}

/*

func TestReadCmd(iotcfg IotConfig) Cmd {

	cmd := Cmd{}
	iotinfo, err := serialTest(iotcfg)
	if err != nil {
		cmd.Para = CMD_dataerror
		cmd.Data = nil

	}
	cmd.Para = CMD_ok
	cmd.Data = iotinfo
	return cmd
}
*/
