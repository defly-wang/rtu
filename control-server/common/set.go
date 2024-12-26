package common

import (
	"fmt"

	"github.com/spf13/viper"
)

type CtlServerSet struct {
	Address string
	Port    int
}

type RtuServerSet struct {
	Address string
	Port    int
}

type WebServerSet struct {
	Address string
	Port    int
}

type RmtServerSet struct {
	Address   string
	Localport []int
	Proxyport []int
}

type Setting struct {
	Ctrserver CtlServerSet
	Rtuserver RtuServerSet
	Webserver WebServerSet
	Rmtserver RmtServerSet
}

var set Setting

func GetRtuAddr() string {
	return fmt.Sprint(set.Rtuserver.Address, ":", set.Rtuserver.Port)
}

func GetCtlAddr() string {
	return fmt.Sprint(set.Ctrserver.Address, ":", set.Ctrserver.Port)
}

func GetWebAddr() string {
	return fmt.Sprint(set.Webserver.Address, ":", set.Webserver.Port)
}

func GetSetting() Setting {
	return set
}

func InitSetting() {

	set = Setting{}

	viper.SetConfigFile("server.toml")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}

	if viper.Unmarshal(&set) != nil {
		fmt.Println(err)
	}
}
