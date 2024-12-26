/*
版本号：2024-6-18
*/
package config

import (
	"fmt"
	"reflect"
	"time"

	"github.com/spf13/viper"
)

//var vipcfg *viper.Viper

type Config struct {
	Common CommonConfig      `toml:"common"`
	Iots   map[int]IotConfig `toml:"iots"`
	Mqtt   MqttConfig        `toml:"mqtt"`
	Webapi WebapiConfig      `toml:"webapi"`
}

type WebapiConfig struct {
	Url   string `json:"url"`
	Token string `json:"token"`
	Used  bool   `json:"used"`
}

type CommonConfig struct {
	Id        string `json:"id"`
	Uint      string `json:"uint"`
	Iotnumber int    `json:"iotnumber"`
	Datadir   string `json:"datadir"`
}

type MqttConfig struct {
	Host      string `json:"host"`
	Port      uint   `json:"port"`
	Pubtoptic string `json:"pubtoptic"`
	Password  string `json:"password"`
	Username  string `json:"username"`
	Used      bool   `json:"used"`
}

type IotConfig struct {
	Iot  int    `toml:"-" json:"iot"`
	Id   string `json:"id"`
	Type string `json:"type"`
	Cron string `json:"cron"`

	Com        string `json:"com"`
	Baudrate   uint   `json:"baudrate"`
	Databits   uint   `json:"databits"`
	Stopbits   uint   `json:"stopbits"`
	Paritymode uint   `json:"paritymode"`

	Read_buff  string `json:"read_buff"`
	Read_len   int    `json:"read_len"`
	Read_delay int    `json:"read_delay"`

	Recive_len      int     `json:"recive_len"`
	Revive_data_len int     `json:"revive_data_len"`
	Revive_data     []int   `json:"revive_data"`
	Offset          bool    `json:"offset"`
	Ratio           float64 `json:"ratio"`
	Crc             string  `json:"crc"`
}

const (
	SYSTEM_RUN_PATH           = "/opt/rtu/"
	SYSTEM_CONFIG_PATH        = "/etc/rtu/"
	SYSTEM_CONFIG_HELPPATH    = "/etc/rtu/help/"
	SYSTEM_CONFIG_IOTHELPPATH = "/etc/rtu/iots/"
	SYSTEM_CONFIG_DATAPATH    = "/mnt/sd/"
	SYSTEM_CONFIG_CGIPATH     = "/usr/boa/www/cgi-bin/"
)
const (
	CONFIG_FILE       = "iotread.toml"
	CONFIG_RTU_FILE   = "rtuclient.toml"
	CONFIG_POWER_FILE = "power.toml"
)

func GetPowerSet() (SerialSet, bool) {
	set := SerialSet{}
	power := viper.New()
	power.SetConfigFile(SYSTEM_CONFIG_PATH + CONFIG_POWER_FILE)
	power.SetConfigType("toml")
	//power.AddConfigPath(".")
	err := power.ReadInConfig()
	if err != nil {
		return set, false
	}

	if power.Unmarshal(&set) != nil {
		return set, false
	}
	return set, true
}

func GetIots(cfg Config) []IotConfig {
	//cfg := GetConfig()
	//fmt.Println(config)

	var result []IotConfig
	for id, iot := range cfg.Iots {
		iot.Iot = id
		result = append(result, iot)
	}
	return result
}

// 查找Iot
func FindIot(cfg Config, iot int) bool {
	_, exist := cfg.Iots[iot]
	return exist
}

func FindIotFromId(cfg Config, id string) (IotConfig, bool) {
	for _, iot := range cfg.Iots {
		if iot.Id == id {
			return iot, true
		}
	}
	return IotConfig{}, false
}

func GetIot(cfg Config, iot int) IotConfig {
	value, exist := cfg.Iots[iot]
	if exist {
		value.Iot = iot
		return value
	} else {
		return IotConfig{}
	}
}

func GetDataPath(cfg Config) string {
	return cfg.Common.Datadir
}

func GetDataDir(cfg Config, t time.Time) string {
	return fmt.Sprintf("%s/%4d/%02d/", GetDataPath(cfg), t.Year(), int(t.Month()))
}

func GetSaveFileName(cfg Config, t time.Time, iot int) string {
	//sid := hex.EncodeToString([]byte(id))
	return fmt.Sprintf("%s/IOT-%02d-%02d.txt", GetDataDir(cfg, t), iot, t.Day())
}

func GetSaveFileNameById(cfg Config, t time.Time, id string) string {
	//sid := hex.EncodeToString([]byte(id))
	info, exist := FindIotFromId(cfg, id)
	if exist {
		return fmt.Sprintf("%s/IOT-%02d-%02d.txt", GetDataDir(cfg, t), info.Iot, t.Day())
	} else {
		return ""
	}
}

func SetIots(oldconfig Config, iots []IotConfig) Config {
	newconfig := oldconfig
	for k := range newconfig.Iots {
		delete(newconfig.Iots, k)
	}
	for _, iot := range iots {
		newconfig.Iots[iot.Iot] = iot
	}

	newconfig.Common.Iotnumber = len(newconfig.Iots)
	return newconfig
	//saveconfig(newconfig)
}

func LoadConfigFile() (Config, error) {
	vipcfg := viper.New()
	filename := SYSTEM_CONFIG_PATH + CONFIG_FILE

	vipcfg.SetConfigFile(filename)
	vipcfg.SetConfigType("toml")

	return toconfig(*vipcfg)

}

func toconfig(vipcfg viper.Viper) (Config, error) {
	result := Config{}

	err := vipcfg.ReadInConfig()
	if err != nil {
		//fmt.Println(err)
		return result, err //return Config{}
	}

	uerr := vipcfg.Unmarshal(&result)
	if uerr != nil {
		//fmt.Println(uerr)
		return result, uerr
	}
	return result, nil
}

/*
func LoadConfig(name string, path []string) (Config, error) {
	vipcfg := viper.New()

	vipcfg.SetConfigName(name)
	vipcfg.SetConfigType("toml")
	for _, p := range path {
		vipcfg.AddConfigPath(p)
	}

	return toconfig(*vipcfg)
	//vipcfg.AddConfigPath("/etc/iotread/")
	//vipcfg.AddConfigPath(".")

}
*/

func SaveConfigFile(newconfig Config) bool {
	cfg := viper.New()

	filename := SYSTEM_CONFIG_PATH + CONFIG_FILE
	//fmt.Println(GetSetting().Configfile)

	cfg.SetConfigFile(filename)
	cfg.SetConfigType("toml")
	//vipcfg.AddConfigPath("/etc/iotread/")
	//cfg.AddConfigPath(".")

	err := cfg.ReadInConfig()
	if err != nil {
		fmt.Println("file:", err)
		return false
	}

	oldconfig := Config{}
	uerr := cfg.Unmarshal(&oldconfig)
	if uerr != nil {
		fmt.Println(err)
		return false
	}
	//newconfig.Common
	//保存数量
	newconfig.Common.Iotnumber = len(newconfig.Iots)
	//CommonConfig.
	if oldconfig.Common != newconfig.Common {
		cfg.Set("common", newconfig.Common)
	}

	if oldconfig.Mqtt != newconfig.Mqtt {
		cfg.Set("mqtt", newconfig.Mqtt)
	}

	for i, iot := range newconfig.Iots {
		oldiot := oldconfig.Iots[i]
		if !reflect.DeepEqual(iot, oldiot) {
			cfg.Set(fmt.Sprint("iots.", i), iot)
		}
	}

	/*
		if oldconfig.Iots != newconfig.Iots {

		}
	*/

	errs := cfg.WriteConfig()
	if errs != nil {
		fmt.Println(errs)
		return false
	}
	return true

}

func LoadIotReadConfig() (Config, bool) {

	//cfg, err := config.LoadConfigFile(GetConfigFileName())
	cfg, err := LoadConfigFile()

	if err != nil {
		return Config{}, false
	}
	return cfg, true

}

func SaveIotReadConfig(newconfig Config) bool {

	return SaveConfigFile(newconfig)

}
