package common

import (
	"config"
	"fmt"

	"github.com/spf13/viper"
)

type Setting struct {
	Mac        string `json:"mac"`
	Unit       string `json:"unit"`
	Id         string `json:"id"`
	Configfile string `json:"-"`
	Server     ServerSet
	// Power      config.SerialSet
}
type ServerSet struct {
	Address     string
	Port        uint
	Signaldelay int
	Watchdog    bool
	Retrydelay  int
}

var set Setting

func GetSetting() Setting {
	return set
}

func GetServerAddress() string {
	return set.Server.Address
}

func GetConfigFileName() string {
	return set.Configfile
}

func GetHSSignalDelay() int {
	return set.Server.Signaldelay
}

func GetConnectRetrylDelay() int {
	return set.Server.Retrydelay
}

/*
func GetPowerSet() config.SerialSet {
	return set.Power
}
*/

func InitSetting() {

	set = Setting{}

	viper.SetConfigFile(config.SYSTEM_CONFIG_PATH + config.CONFIG_RTU_FILE)
	viper.SetConfigType("toml")
	//viper.AddConfigPath(".")
	//viper.AddConfigPath(config.SYSTEM_CONFIG_PATH)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}

	if viper.Unmarshal(&set) != nil {
		fmt.Println(err)
	}
}
