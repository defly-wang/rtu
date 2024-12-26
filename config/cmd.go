package config

import (
	"errors"
	"strings"

	"github.com/spf13/viper"
)

// cmd
const (
	//Control发起  Server 直接返回数据
	//src=CTLip，dec="";返回 src="",dec=CTLip,data={Clients=RTU} para=CMD_PARA_GET
	CMD_RTUClients = "rtu"
	//心跳信号 src=RTU ip ，dec=""，date=MAC地址；返回：src=""，dec=RTUip data="ok"
	CMD_HeartbeatSignal = "hs"

	//获取RTU系统信息 src=CTLip,dec=RTUip;返回：src=RTUip，dec=CTLip，data={Sysinfo数据} para=CMD_PARA_GET
	CMD_Info = "info"
	//获取RTU运行状态 src=CTLip,dec=RTUip;返回：src=RTUip，dec=CTLip，data={Sysinfo数据} para=CMD_PARA_GET
	CMD_Status = "status"

	//获取/设置基本设置 para=CMD_PARA_GET，edit
	CMD_Base = "base"
	//获取/设置Mqtt设置 para=CMD_PARA_GET，edit
	CMD_Mqtt = "mqtt"
	//获取/设置Webapi设置 para=CMD_PARA_GET，edit
	CMD_Webapi = "webapi"
	//获取/设置 RTU 传感器设置（所有）para=CMD_PARA_GET=all/iot，edit=iot,del=iot
	//参数值传入 "get=all" “get=雨量计1” "del=1" ...
	//read=1,读取iot1数据
	CMD_Iot = "iot"

	//测试设备
	//para=read data=IotConfig
	//para=mqtt date=MqttConfig
	//para=webapi data=WebapiConfig
	CMD_Test = "test"

	//自动生成IOT内容（根据设备类型、地址等生成读取命令（CRC校验）、解析字段等
	//para=get data=SignalRead 返回读取命令（hexString格式）
	//para=read data=厂商-类型（字符串）  //读取 厂商-类型.toml文件，返回IotConfig
	//para=mqtt 返回缺省设置
	//para=webapi 返回缺省设置
	CMD_Help = "help"

	//历史数据
	//请求
	//"para"="day=2024-12-12=1"  一天数据
	//"para"="one=2024-12-12 08:00:00=2"  一个数据，距离时间最近
	//"para"="count=1"
	//返回 data={}
	CMD_Histiry = "his"

	//远控frpc 远控某电脑
	//para=get ，dec=目标，返回列表（port)
	//{"cmd":"remote","para":"get","dec":"192.168.1.129:34122"}

	//返回值中：data "Localport":6001 是本地访问端口，将自动proxy到 "Remoteport":80

	//{"cmd":"remote","src":"192.168.1.129:34122","type":"","dec":"","para":"ok",
	//"data":[{"Ip":"192.168.1.129:34122","Remoteport":80,"Localport":6001,"Locallistener":{},"Clients":null}]}

	//para=add, dec=目标 data="port" 设置端口号 字符串类型 "data":"22"
	//{"cmd":"remote","para":"add","dec":"192.168.1.129:34122","data":"22"}
	//para=del, dec=目标
	//{"cmd":"remote","para":"del","dec":"192.168.1.129:34122"}
	CMD_Remote = "remote"

	//电源信息
	//para="info" 电源信息
	//para="data" 运行数据
	//para="set"  设置数据

	//para=edit 设置
	//data= 设置电池类型：//1-4,//电池类型  0=自定义， 1=开口， 2=密封， 3=胶体， 4=锂电池；容量AH
	//type PowerTypeCapacity struct {
	//	BatterCapacity uint16 `json:"battercapacity"`
	//	BatterType     uint16 `json:"battertype"`
	//}

	CMD_Power = "power"

	//CMD_Read = "readiot"

	//三种连接客户，一种是RTU设备，二是控制客户端（socket），三是websocket控制段
	ClientTypeRTU = "RTU"
	ClientTypeCTL = "CTL"
	ClientTypeWEB = "WEB"

	CMD_PARA_GET  = "get"
	CMD_PARA_EDIT = "edit"
	CMD_PARA_DEL  = "del"
	CMD_PARA_ADD  = "add"

	CMD_PARA_INFO = "info"
	CMD_PARA_DATA = "data"
	CMD_PARA_SET  = "set"

	CMD_PARA_READ   = "read"
	CMD_PARA_MQTT   = "mqtt"
	CMD_PARA_WEBAPI = "webapi"

	CMD_PARA_HIS_ONE = "one"
	CMD_PARA_HIS_DAY = "day"

	CMD_PARA_ALL = "all"

	//返回值 一般写入para参数中
	CMD_ok    = "ok"
	CMD_error = "error"

	CMD_paraerror = "paraerror"
	CMD_dataerror = "dataerror"
	CMD_noexist   = "noexist"
	CMD_existed   = "existed"
	CMD_timeout   = "timeout"
	CMD_unsupport = "unsupported"
)

//命令格式，发起命令、返回命令
/*
cmd=命令
para=参数
get  获取数据
edit 修改数据
add  添加数据（iot）
del  删除数据（iot）



src=发起方，存储发起方IP与端口  使用conn.RemoteAddr().String()获取
type=发起方类型 主要区别ctl（socket） web（websocket）
dec=发送目的
*/

type Cmd struct {
	Cmd  string `json:"cmd"`
	Src  string `json:"src"`
	Type string `json:"type"`
	Dec  string `json:"dec"`
	Para string `json:"para"`
	Data any    `json:"data"`
	/*
	   Mac  string `json:"mac"`
	*/
}

func GetDataFromCmd(buff []byte, data any) error {
	dc := viper.New()
	dc.SetConfigType("json")

	err := dc.ReadConfig(strings.NewReader(string(buff)))
	if err != nil {
		//fmt.Println(err)
		return err
	}

	if dc == nil {
		return errors.New("data nil err")
	}

	sub := dc.Sub("data")
	if sub == nil {
		return errors.New("data nil err")
	}

	erru := sub.Unmarshal(&data)
	if erru != nil {
		//fmt.Println(erru)
		return erru
	}
	return nil
}
