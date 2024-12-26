package common

import (
	"github.com/spf13/viper"
)

const (
	CfgFile = "cgi.toml"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetUserInfo() (User, bool) {

	user := User{}

	cfgvip := viper.New()
	cfgvip.SetConfigFile(CfgFile)
	cfgvip.SetConfigType("toml")
	err := cfgvip.ReadInConfig()
	if err != nil {
		//fmt.Println(err)
		return user, false
	}

	errm := cfgvip.Unmarshal(&user)
	if errm != nil {
		return user, false
	}
	/*
		else {
			user.Username = viper.GetString("common.username")
			user.Password = viper.GetString("common.password")
		}
	*/

	return user, true

}

func SavePassword(pwd string) bool {

	cfgvip := viper.New()
	cfgvip.SetConfigFile(CfgFile)
	cfgvip.SetConfigType("toml")

	if cfgvip.ReadInConfig() != nil {
		return false
	}

	cfgvip.Set("password", pwd)

	return cfgvip.WriteConfig() == nil
}

/*
func GetConfigFileName() string {
	viper.SetConfigFile("set.toml")
	err := viper.ReadInConfig()
	if err != nil {
		return ""
	}
	return viper.GetString("common.configfile")
}
*/
