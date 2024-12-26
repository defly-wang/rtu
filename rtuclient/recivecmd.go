package main

import (
	"config"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ProcRecive(recive []byte) {

	cmd := config.Cmd{}
	jerr := json.Unmarshal(recive, &cmd)
	if jerr != nil {
		fmt.Println(jerr.Error())
		return
	}
	sendcmd := config.Cmd{
		Cmd:  cmd.Cmd,
		Dec:  cmd.Src,
		Src:  cmd.Dec,
		Type: cmd.Type,
		Para: config.CMD_ok,
	}
	//fmt.Println(cmd)
	switch cmd.Cmd {
	case config.CMD_HeartbeatSignal:
		//do nothing
		//写watch dog
		go config.LedBlink2(800)

		hstime = time.Now()
		//fmt.Println("hk")

		return

	case config.CMD_Remote:

		proxyinfo := config.Proxy{}
		config.GetDataFromCmd(recive, &proxyinfo)
		//fmt.Println(proxyinfo)
		ok, ip := StartRemote(proxyinfo.Remoteport, proxyinfo.Proxyport)
		if ok {
			proxyinfo.Ip = ip
			sendcmd.Data = proxyinfo
		} else {
			sendcmd.Data = cmd.Data
			sendcmd.Para = config.CMD_error
		}

		//SendData(sendcmd)
		//return

	case config.CMD_Info:
		//fmt.Println("iii cmd=", cmd)
		//可去掉参数
		if cmd.Para == config.CMD_PARA_GET {
			info := config.GetSysInfo()
			//info.Iotread_ver = config.GetIotVersion()
			//info.Rtuclient_ver = VERSION
			//info.Webcgi_ver = config.GetCgiVersion()

			sendcmd.Data = info
		} else {
			sendcmd.Para = config.CMD_paraerror
		}

	case config.CMD_Status:
		if cmd.Para == config.CMD_PARA_GET {
			sendcmd.Data = config.GetRtuRunStatus()
		} else {
			sendcmd.Para = config.CMD_paraerror
		}

	case config.CMD_Iot:
		para := strings.Split(cmd.Para, "=")
		if len(para) < 2 {
			sendcmd.Para = config.CMD_paraerror
			break
		}
		CMDIots(&sendcmd, para[0], para[1], recive)

	case config.CMD_Mqtt:
		switch cmd.Para {
		case config.CMD_PARA_GET:
			cfg, ok := config.LoadIotReadConfig()
			if ok {
				sendcmd.Data = cfg.Mqtt
			} else {
				sendcmd.Para = config.CMD_dataerror
			}

		case config.CMD_PARA_EDIT:
			mqtt := config.MqttConfig{}
			err := config.GetDataFromCmd(recive, &mqtt)
			if err == nil {
				newconfig, ok := config.LoadIotReadConfig()
				if ok {
					newconfig.Mqtt = mqtt
					if !config.SaveIotReadConfig(newconfig) {
						sendcmd.Para = config.CMD_dataerror
					}
				} else {
					sendcmd.Para = config.CMD_dataerror
				}
			} else {
				sendcmd.Para = config.CMD_dataerror
			}
			//common.SaveMqtt(mqtt)
		default:
			sendcmd.Para = config.CMD_paraerror
		}

	case config.CMD_Base:
		switch cmd.Para {
		case config.CMD_PARA_GET:
			cfg, ok := config.LoadIotReadConfig()
			if ok {
				sendcmd.Data = cfg.Common
			} else {
				sendcmd.Para = config.CMD_paraerror
			}
		case config.CMD_PARA_EDIT:
			base := config.CommonConfig{}
			err := config.GetDataFromCmd(recive, &base)
			if err == nil {
				newconfig, ok := config.LoadIotReadConfig()
				if ok {
					newconfig.Common = base
					//common.SaveBase(base)
					if !config.SaveIotReadConfig(newconfig) {
						sendcmd.Para = config.CMD_dataerror
					}
				} else {
					sendcmd.Para = config.CMD_dataerror
				}

			} else {
				sendcmd.Para = config.CMD_dataerror
			}
		default:
			sendcmd.Para = config.CMD_paraerror
		}

	case config.CMD_Webapi:
		switch cmd.Para {
		case config.CMD_PARA_GET:
			cfg, ok := config.LoadIotReadConfig()
			if ok {
				sendcmd.Data = cfg.Webapi
			} else {
				sendcmd.Para = config.CMD_dataerror
			}
		case config.CMD_PARA_EDIT:
			api := config.WebapiConfig{}
			err := config.GetDataFromCmd(recive, &api)
			if err == nil {
				newconfig, ok := config.LoadIotReadConfig()
				if ok {
					newconfig.Webapi = api
					if !config.SaveIotReadConfig(newconfig) {
						sendcmd.Para = config.CMD_dataerror
					}
				} else {
					sendcmd.Para = config.CMD_dataerror
				}
			} else {
				sendcmd.Para = config.CMD_dataerror
			}
			//common.SaveWebapi(api)
		default:
			sendcmd.Para = config.CMD_paraerror
		}

	case config.CMD_Histiry:
		para := strings.Split(cmd.Para, "=")
		if len(para) < 3 {
			sendcmd.Para = config.CMD_paraerror
			break
		}

		t, err := config.StrToTime(para[1])
		//	t, err := time.Parse(format, para[1])
		if err != nil {
			sendcmd.Para = config.CMD_paraerror
			break
		}

		iot, erriot := strconv.Atoi(para[2])
		if erriot != nil {
			sendcmd.Para = config.CMD_paraerror
			break
		}

		//处理可能错误
		switch para[0] {
		case config.CMD_PARA_HIS_ONE:
			//t, _ := time.Parse("2006-01-02 15:04:05", para[1])
			//iot, _ := strconv.Atoi(para[2])
			cfg, ok := config.LoadIotReadConfig()
			if ok {
				data, errl := config.LoadHistoryDataOne(cfg, t, iot)
				if errl != nil {
					break
				}
				sendcmd.Data = data
			} else {
				sendcmd.Para = config.CMD_paraerror
				break
			}
		case config.CMD_PARA_HIS_DAY:
			//t, _ := time.Parse("2006-01-02", para[1])
			//iot, _ := strconv.Atoi(para[2])
			cfg, ok := config.LoadIotReadConfig()
			if !ok {
				sendcmd.Para = config.CMD_paraerror
				break
			}
			datas, errl := config.LoadHistoryData(cfg, t, iot)
			if errl != nil {
				break
			}
			//sendcmd.Data = datas
			//datasize := len(datas)
			var tmpdate []config.IotInfo

			j := 0
			for _, data := range datas {

				tmpdate = append(tmpdate, data)
				j++
				if j > 6 {
					j = 0
					sendcmd.Data = tmpdate
					SendData(sendcmd)
					tmpdate = tmpdate[:0]
				}
			}
			if len(tmpdate) > 0 {
				sendcmd.Data = tmpdate
				SendData(sendcmd)
			}
			return

		default:
		}

	case config.CMD_Power:
		/*
			para := strings.Split(cmd.Para, "=")
			if len(para) < 2 {
				sendcmd.Para = config.CMD_paraerror
				break
			}
			if para[0] != config.CMD_PARA_GET {
				sendcmd.Para = config.CMD_paraerror
				break
			}
		*/
		CMDPower(&sendcmd, cmd.Para, recive)

	case config.CMD_Help:
		switch cmd.Para {

		case config.CMD_PARA_GET:
			read := config.SignalRead{}
			err := config.GetDataFromCmd(recive, &read)
			if err != nil {
				sendcmd.Para = config.CMD_paraerror
				break
			}
			sendcmd.Data = hex.EncodeToString(config.HelpSerialReadCmd(read))
		case config.CMD_PARA_READ:
			name, ok := cmd.Data.(string)
			if ok {
				iotinfo, err := config.GetIotHelpConfig(name)
				if err != nil {
					sendcmd.Para = config.CMD_paraerror
					sendcmd.Data = nil
				} else {
					sendcmd.Data = iotinfo
				}
			} else {
				sendcmd.Para = config.CMD_paraerror
			}
		case "org":
			sendcmd.Data = config.LoadHelpIotOrg()
		case "type":
			sendcmd.Data = config.LoadHelpIotType()
		case config.CMD_PARA_MQTT:
			mqtt, ok := config.LoadHelpMqtt()
			if ok {
				sendcmd.Data = mqtt
			} else {
				sendcmd.Para = config.CMD_dataerror
			}

		case config.CMD_PARA_WEBAPI:
			webapi, ok := config.LoadHelpWebapi()
			if ok {
				sendcmd.Data = webapi
			} else {
				sendcmd.Para = config.CMD_dataerror
			}
		default:
			sendcmd.Para = config.CMD_paraerror
		}

	case config.CMD_Test:
		switch cmd.Para {
		case config.CMD_PARA_READ:
			iotcfg := config.IotConfig{}
			err := config.GetDataFromCmd(recive, &iotcfg)
			if err != nil {
				sendcmd.Para = config.CMD_paraerror
				break
			}
			iotinfo, errr := config.TestReadSirial(iotcfg)
			if errr != nil {
				sendcmd.Para = config.CMD_dataerror
				sendcmd.Data = errr.Error()
				break
			}
			sendcmd.Data = iotinfo
		case config.CMD_PARA_MQTT:
			mqttinfo := config.MqttConfig{}
			err := config.GetDataFromCmd(recive, &mqttinfo)
			if err != nil {
				sendcmd.Para = config.CMD_dataerror
				sendcmd.Data = "MQTT参数错误！"
				break
			}
			errm := config.TestMqtt(mqttinfo)
			if errm != nil {
				sendcmd.Para = config.CMD_dataerror
				sendcmd.Data = errm.Error()
				break
			}
			sendcmd.Data = config.CMD_ok

		case config.CMD_PARA_WEBAPI:
			api := config.WebapiConfig{}
			err := config.GetDataFromCmd(recive, &api)
			if err != nil {
				sendcmd.Para = config.CMD_dataerror
				sendcmd.Data = "Webapi参数错误！"
				break
			}
			errt := config.TestWebapi(api)
			if errt != nil {
				sendcmd.Para = config.CMD_dataerror
				sendcmd.Data = errt.Error()
				break
			}
			sendcmd.Data = config.CMD_ok

			//break

		default:
			break
		}

	default:
		sendcmd.Para = config.CMD_unsupport // "unsupported"
	}

	SendData(sendcmd)
}

//func CMDPower

func CMDPower(sendcmd *config.Cmd, para string, recive []byte) {
	power, ok := config.GetPowerSet()
	if !ok {
		sendcmd.Para = config.CMD_dataerror
		return
	}

	switch para {
	case config.CMD_PARA_INFO:
		powerinfo, err := config.ReadPowerInfo(power)
		if err != nil {
			sendcmd.Para = config.CMD_dataerror
			sendcmd.Data = nil
		} else {
			sendcmd.Data = powerinfo
		}
	case config.CMD_PARA_DATA:
		powerdata, err := config.ReadPowerRunData(power)
		if err != nil {
			sendcmd.Para = config.CMD_dataerror
			sendcmd.Data = nil
		} else {
			sendcmd.Data = powerdata
		}
	case config.CMD_PARA_SET:
		powerset, err := config.ReadPowerSet(power)
		if err != nil {
			sendcmd.Para = config.CMD_dataerror
			sendcmd.Data = nil
		} else {
			sendcmd.Data = powerset
		}
		/*
			case config.CMD_PARA_EDIT:
				err := config.SetPowerBatterType(common.GetPowerSet(), 4)
				if err != nil {
					sendcmd.Para = config.CMD_dataerror
					sendcmd.Data = nil
				} else {
					sendcmd.Data = "ok"
				}
		*/
	case config.CMD_PARA_EDIT:
		ptc := config.PowerTypeCapacity{}
		err := config.GetDataFromCmd(recive, &ptc)
		if err != nil {
			sendcmd.Para = config.CMD_dataerror
			sendcmd.Data = "输入信息错误！"
			break
		}
		if config.SetPowerBatterType(power, ptc.BatterType) != nil {
			sendcmd.Para = config.CMD_dataerror
			sendcmd.Data = "写入设置信息错误！"
			break
		}
		if config.SetPowerBatterCapacity(power, ptc.BatterCapacity) != nil {
			sendcmd.Para = config.CMD_dataerror
			sendcmd.Data = "写入设置信息错误！"
			break
		}
		sendcmd.Data = "ok"

	default:
		sendcmd.Para = config.CMD_paraerror
	}
}

func CMDIots(sendcmd *config.Cmd, cmd, para string, recive []byte) {
	cfg, ok := config.LoadIotReadConfig()
	if !ok {
		sendcmd.Para = config.CMD_dataerror
		return
	}
	switch cmd {
	case config.CMD_PARA_GET:

		if para == config.CMD_PARA_ALL {
			sendcmd.Data = config.GetIots(cfg)
		} else {

			iot, err := strconv.Atoi(para)
			if err == nil {
				if config.FindIot(cfg, iot) {
					sendcmd.Data = config.GetIot(cfg, iot)
				} else {
					sendcmd.Para = config.CMD_noexist
				}
			} else {
				sendcmd.Para = config.CMD_paraerror
			}
		}

	case config.CMD_PARA_ADD:

		iot, err := strconv.Atoi(para)
		if err == nil {
			iotdata := config.IotConfig{}
			errd := config.GetDataFromCmd(recive, &iotdata)
			if errd == nil {
				if !config.FindIot(cfg, iot) {
					//cfg := common.GetConfig()
					cfg.Iots[iot] = iotdata

					if !config.SaveIotReadConfig(cfg) {
						sendcmd.Para = config.CMD_dataerror
					}
				} else {
					sendcmd.Para = config.CMD_existed
				}
			} else {
				sendcmd.Para = config.CMD_dataerror
			}
		} else {
			sendcmd.Para = config.CMD_paraerror
		}

	case config.CMD_PARA_DEL:
		iot, err := strconv.Atoi(para)
		if err == nil {
			if config.FindIot(cfg, iot) {
				//cfg := common.GetConfig()
				delete(cfg.Iots, iot)
				if !config.SaveIotReadConfig(cfg) {
					sendcmd.Para = config.CMD_dataerror
				}
			} else {
				sendcmd.Para = config.CMD_noexist
			}
		} else {
			sendcmd.Para = config.CMD_paraerror
		}
	case config.CMD_PARA_EDIT:
		iot, err := strconv.Atoi(para)
		if err == nil {
			iotdata := config.IotConfig{}
			errd := config.GetDataFromCmd(recive, &iotdata)
			if errd == nil {
				if config.FindIot(cfg, iot) {
					//cfg := common.GetConfig()
					cfg.Iots[iot] = iotdata
					if !config.SaveIotReadConfig(cfg) {
						sendcmd.Para = config.CMD_dataerror
					}
				} else {
					sendcmd.Para = config.CMD_noexist
				}
			} else {
				sendcmd.Para = config.CMD_dataerror
			}
		} else {
			sendcmd.Para = config.CMD_paraerror
		}

	case config.CMD_PARA_READ:
		iot, err := strconv.Atoi(para)
		if err == nil {

			iotcfg := cfg.Iots[iot]

			d, errs := config.TestReadSirial(iotcfg)
			if errs == nil {
				sendcmd.Data = d
			} else {
				sendcmd.Para = config.CMD_dataerror
			}

		} else {
			sendcmd.Para = config.CMD_dataerror
		}

	default:
		sendcmd.Para = config.CMD_paraerror
	}

}
