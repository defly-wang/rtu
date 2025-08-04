package send

import (
	"fmt"
	"iotread/common"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var mtqqClient mqtt.Client

func MqttInit() bool {

	opts := mqtt.NewClientOptions()
	server := fmt.Sprint("tcp://", common.CFG.Mqtt.Host, ":", common.CFG.Mqtt.Port)
	//fmt.Println(server)
	opts.AddBroker(server)
	//append(opts.Servers, )
	opts.SetClientID(common.MakeMqttClientID())
	opts.SetUsername(common.CFG.Mqtt.Username) //config.GetMqttUser())
	opts.SetPassword(common.CFG.Mqtt.Password)

	mtqqClient = mqtt.NewClient(opts)

	if token := mtqqClient.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return false
	}
	return true
}

func MqttClose() {
	mtqqClient.Disconnect(300)
}

func MqttPublicInfo(toptic string, jsondata string) bool {

	if !mtqqClient.IsConnected() {
		//断线重联
		if !MqttInit() {
			return false
		}
		//return false
	}
	//fmt.Println(string(jsondata))

	if token := mtqqClient.Publish(toptic, 0, true, jsondata); token.Wait() && token.Error() != nil {
		//fmt.Println(token.Error())
		return false
	}
	return true
}
