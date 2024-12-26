/*
版本号：2024-6-18
*/
package common

import (
	"config"
	"encoding/hex"
	"time"
)

var VERSION = "1.0.0"

/*
var (
	name          = "config"
	path []string = []string{".", "/etc/rtu/"}
)
*/

var CFG config.Config

func InitConfig() bool {
	var err error
	//CFG, err = config.LoadConfig(name, path) // LoadConfig(name, path)
	CFG, err = config.LoadConfigFile()
	return err == nil
}

func InitIotInfo(iot int) config.IotInfo {
	return config.IotInfo{
		Iot: iot,

		Id:    CFG.Iots[iot].Id,   // GetIotid(iot),
		Type:  CFG.Iots[iot].Type, //tGetIotType(iot),
		Time:  config.MyTime(time.Now()),
		Ratio: CFG.Iots[iot].Ratio, //GetIotRatio(iot),
	}
}

func GetIotPortName(iot int) string {
	//fmt.Println("/dev/ttyS" + CFG.Iots[iot].Com)
	return "/dev/ttyS" + CFG.Iots[iot].Com
}

func GetIotReadbuff(iot int) []byte {
	sendbuff, err := hex.DecodeString(CFG.Iots[iot].Read_buff)
	if err != nil {
		//fmt.Println(err)
		return []byte{0}
	}
	return sendbuff
}

func MakeMqttClientID() string {

	//iotnum := viper.GetInt("common.iotnumber")
	var iotids = "("

	for iot := range CFG.Iots {
		iotids += CFG.Iots[iot].Id
		if iot < CFG.Common.Iotnumber {
			iotids += ","
		} else {
			iotids += ")"
		}
	}

	//commond.uint-commond.id(iot1.id,iot2.id)
	return CFG.Common.Uint + "-" + iotids
	//return viper.GetString("common.uint") + "-" + viper.GetString("common.id") + iotids
}
