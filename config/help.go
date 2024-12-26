package config

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/spf13/viper"
)

// const SYSTEM_CONFIG_HELP = "/etc/rtu/iothelp.toml"
const (
	//SYSTEM_CONFIG_HELP = "/etc/rtu/iothelp.toml"
	CONFIG_FILE_IOTHELP    = "iothelp.toml"
	CONFIG_FILE_MQTTHELP   = "mqtt.toml"
	CONFIG_FILE_WEBAPIHELP = "webapi.toml"
	//SYSTEM_CONFIG_HELP        = "./iothelp.toml"
	//SYSTEM_CONFIG_HELP_MQTT   = "./mqtt.toml"
	//SYSTEM_CONFIG_HELP_WEBAPI = "./webapi.toml"
)

type SignalRead struct {
	Add     uint8  `json:"add"`
	Cmd     uint8  `json:"cmd"`
	DataAdd uint16 `json:"dataadd"`
	DataLen uint16 `json:"datalen"`
}

type IotType struct {
	Org     string
	OrgName string
	Type    string
	Name    string
	Mode    string
	File    string
}

type IotOrg struct {
	Org  string
	Name string
	//Type    []IotType
	//ConfigFile string
}

func LoadHelpIotOrg() []IotOrg {
	vipcfg := viper.New()

	cfgfile := SYSTEM_CONFIG_HELPPATH + CONFIG_FILE_IOTHELP

	vipcfg.SetConfigFile(cfgfile)
	vipcfg.SetConfigType("toml")

	err := vipcfg.ReadInConfig()
	if err != nil {
		return nil
	}

	var orgs []IotOrg
	for _, org := range vipcfg.Get("org").([]interface{}) {
		var o IotOrg
		e := mapToStruct(org.(map[string]interface{}), &o)
		if e == nil {
			orgs = append(orgs, o)
		}
	}
	return orgs
}

func LoadHelpIotType() []IotType {
	vipcfg := viper.New()

	cfgfile := SYSTEM_CONFIG_HELPPATH + CONFIG_FILE_IOTHELP
	vipcfg.SetConfigFile(cfgfile)
	vipcfg.SetConfigType("toml")

	err := vipcfg.ReadInConfig()
	if err != nil {
		return nil
	}

	orgs := LoadHelpIotOrg()

	//fmt.Println(orgs)
	//orgs := vipcfg.Get("org")

	var iothelps []IotType
	//helps :=

	for _, help := range vipcfg.Get("iothelp").([]interface{}) {
		var iothelp IotType
		e := mapToStruct(help.(map[string]interface{}), &iothelp)
		if e == nil {
			iothelp.OrgName = getOrgName(orgs, iothelp.Org)
			iothelps = append(iothelps, iothelp)
		}
	}
	//fmt.Println(iothelps)
	return iothelps
}

func LoadHelpMqtt() (MqttConfig, bool) {
	cfg := viper.New()

	cfgfile := SYSTEM_CONFIG_HELPPATH + CONFIG_FILE_MQTTHELP

	cfg.SetConfigFile(cfgfile)
	cfg.SetConfigType("toml")

	cfg.ReadInConfig()

	mqtt := MqttConfig{}
	err := cfg.Unmarshal(&mqtt)
	if err != nil {
		return mqtt, false
	}
	return mqtt, true
}

func LoadHelpWebapi() (WebapiConfig, bool) {
	cfg := viper.New()

	cfgfile := SYSTEM_CONFIG_HELPPATH + CONFIG_FILE_WEBAPIHELP

	cfg.SetConfigFile(cfgfile)
	cfg.SetConfigType("toml")

	cfg.ReadInConfig()

	webapi := WebapiConfig{}
	err := cfg.Unmarshal(&webapi)
	if err != nil {
		return webapi, false
	}
	return webapi, true
}

func getOrgName(orgs []IotOrg, org string) string {
	for _, o := range orgs {
		if o.Org == org {
			return o.Name
		}
	}
	return ""
}

func HelpSerialReadCmd(sr SignalRead) []byte {
	buff := make([]byte, 8)

	buff[0] = sr.Add
	buff[1] = sr.Cmd
	buff[2] = uint8((sr.DataAdd & 0xff00) >> 8)
	buff[3] = uint8((sr.DataAdd & 0x00ff))
	buff[4] = uint8((sr.DataLen & 0xff00) >> 8)
	buff[5] = uint8((sr.DataLen & 0x00ff))
	crc16 := makeCRC16(buff[:6])
	buff[6] = byte(crc16 & 0xff)
	buff[7] = byte(crc16 >> 8)
	return buff
}

func GetIotHelpConfig(filename string) (IotConfig, error) {
	result := IotConfig{}

	toml := ".toml"
	if strings.Contains(filename, ".toml") {
		toml = ""
	}

	cfgfile := SYSTEM_CONFIG_IOTHELPPATH + filename + toml
	iotcfg := viper.New()

	iotcfg.SetConfigFile(cfgfile)
	iotcfg.SetConfigType("toml")

	err := iotcfg.ReadInConfig()
	if err != nil {
		return result, err //return Config{}
	}

	uerr := iotcfg.Unmarshal(&result)
	if uerr != nil {
		//fmt.Println(uerr)
		return result, uerr
	}
	return result, nil

}

func mapToStruct(m map[string]interface{}, s interface{}) error {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr && !v.Elem().CanSet() {
		return fmt.Errorf("struct must be a pointer")
	}

	u := v.Elem()

	for k, val := range m {
		fName := capitalizeFirstLetter(k)

		//fmt.Println(fName)

		//u.FieldByName()
		f := u.FieldByName(fName)

		if !f.IsValid() {

			continue // ignore unexported fields
		}
		if !f.CanSet() {

			continue // ignore fields that are unsettable
		}

		inVal := reflect.ValueOf(val)
		if inVal.Type().AssignableTo(f.Type()) {
			f.Set(inVal)
		} else {
			fmt.Printf("Field: %s, Value type doesn't match\n", fName)
		}
	}
	return nil
}

func capitalizeFirstLetter(str string) string {
	if len(str) == 0 {
		return ""
	}
	runes := []rune(str)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
