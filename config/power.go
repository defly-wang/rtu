package config

import (
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/jacobsa/go-serial/serial"
)

type PowerInfo struct {
	//最大支持电压、额定充电电流、额定放电电流、产品类型、型号、软件版本、硬件版本、序列号
	//0A                    0B                         0C -13         |14-15        |16-17      |18-19   |20   |
	//|01        |02        |03        |04             | 05-20（8字）  |21-24        |25-28      |29-32   |33-34|
	//最大支持电压 |额定充电电流|额定放电电流|产品类型         |型号          |软件版本      |硬件版本     |序列号   |地址  |
	//10进制+V    |10进制+A   |10进制+A  |00控制器、02逆变器|（ascii码     |后三字节每字节代表一个版本子号|转成hex串|低字节 |
	MaxSuppoptVoltage    string
	RateChargeCurrent    string
	RateDischargeCurrent string
	ProductType          string
	ProductModel         string
	SoftwareVer          string
	HardwareVer          string
	SN                   string
}

type PowerRunData struct {
	//电池电量%，蓄电池电压V0.1，充电电流A0.01,控制器温度（b7符号，b0-6值），电池温度，负载直流电压，负载直流电流A0.01，负载功率W
	//0100H------------------------------70字节-----------------------------------------------------------0120H
	//0100-------------------------------14字节--------------------------------------------------------------0106
	//|01-02             |03-04         |05-06      |07                      |08     |09-10        |11-12           |13-14   |
	//|电池电量02（0-100）%|蓄电池电压V 0.1|充电电流A0.01|控制器温度（b7符号，b0-6值）|电池温度|负载直流电压V0.1|负载直流电流A0.01|负载功率W|
	BatterySOC     uint
	BatteryVoltage MyFloat
	ChargeCurrent  MyFloat
	TempControl    int
	TempBattery    int
	LoadVoltage    MyFloat
	LoadCurrent    MyFloat
	LoadPower      uint16

	//太阳能板电压，太阳能板电流，充电功率
	//0107-------------------------------太阳能板信息（6 字节）--------------------------0109H
	//|01-02     |03-04     |05-06  |
	//|太阳能板电压|太阳能板电流|充电功率|
	SolarVoltage MyFloat
	SolarCurrent MyFloat
	SolarPower   uint16

	//开/关负载命令,蓄电池当天最低电压,蓄电池当天最高电压,当天充电最大电流,当天放电最大电流,当天充电最大功率,当天放电最大功率,当天充电安时数,当天放电安时数,当天发电量,当天用电量
	//010AH-----------------------------蓄电池信息（22字节）-----------------------------0114H
	//|01-02                       |03-04            |05-06  （电压*0.1，电流*0.01，功率原值）
	//|开/关负载命令0000关闭，0001打开|蓄电池当天最低电压0.1|
	LoadSwitch          bool
	BatterVoltageMin    MyFloat
	BatterVoltageMax    MyFloat
	ChargeCurrentMax    MyFloat
	DischargeCurrentMax MyFloat
	ChargePowerMax      uint16
	DischargePowerMax   uint16
	ChargeAH            uint16
	DischargeAH         uint16
	GenPower            uint16
	UsePower            uint16

	//0115H-----------------------------历史数据信息（22字节）--------------------------011FH
	//总运行天数，蓄电池总过放次数，蓄电池总充满次数，蓄电池总充电安时数（4b），蓄电池总放电安时数（4b)，累计发电量（4b），累计用电量（4b）
	RunDaysCount       uint16
	OverDischargeCount uint16
	FillingsCount      uint16
	ChargeAHCount      uint32
	DischargeAHCount   uint32
	GenPowerCount      uint32
	UsePowerCount      uint32

	//0120
	//负载状态，负载亮度，充电状态
	//充电状态：00H:未开启充电 01H:启动充电模式 02 mppt充电模式  03H:均衡充电模式 04H:提升充电模式 05H:浮充充电模式 06H:限流(超功率)
	//----------------------------------------负载信息（2字节）------------------------

	LoadStatus   bool
	LoadRright   int
	ChargeStatus int

	//0121-0122
	//故障码总计32位代表21种错误，其他保留
	FaultCode    uint32
	FaultStrings []string
}

type PowerSet struct {
	BatterCapacity        uint16  //蓄电池标称容量  E002H
	VoltageSet            uint8   //系统电压设置    E003H
	VoltageIden           uint8   //识别后的电压
	BatterType            uint16  //电池类型       E004H 0=自定义， 1=开口， 2=密封， 3=胶体， 4=锂电池
	VoltageOver           MyFloat //超压电压    E005H *0.1
	VoltageChargeLim      MyFloat //充电限制电压 E006H *0.1
	VoltageChargeEqu      MyFloat //均衡充电电压 E007H *0.1
	VoltageChargeUpg      MyFloat //提升充电电压 E008H *0.1
	VoltageChargeFlo      MyFloat //浮充充电电压 E009H *0.1
	VoltageChargeUpgR     MyFloat //提升充电返回电压 E00AH *0.1
	VoltageDischargeOverR MyFloat //过放返回电压 E00BH *0.1
	VoltageUnderWarn      MyFloat //欠压警告电压 E00CH *0.1
	VoltageDischargeOver  MyFloat //过放电压    E00DH *0.1
	VoltageDisChargeLim   MyFloat //放电限制电压 E00EH *0.1
	SocChargeLim          uint8   //充电截止SOC E00FH  未实现
	SocDischargeLim       uint8   //放电截止SOC E00FH
	DelaySDischargeOver   uint16  //过放延迟时间（S） E010H   0~120
	TimeMChargeEqu        uint16  //均衡充电时间（M)  E011H   0～300 步长 10
	TimeMChargeUpg        uint16  //提升充电时间（M)  E012H   10～300 步长 10
	IntervalDChargeUpg    uint16  //均衡充电时间（D)  E013H   0～255 0：关闭，步长 5
	TempComp              uint16  //温度补偿系数（mV/C/2V）E014H   0～6 0：不补偿，步长 1
	LoadMode              uint16  //负载工作模式      E01DH   00H 纯光控、光控开/关负载 01H-0EH 光控开负载,延时1-14小时后关闭 0FH(十进制 15) 手动模式 10H(十进制16) 调试模式 11H(十进制17) 常开模式
	DelayMLightCtl        uint16  //光控延时时间（Min）E01EH  0-60
	VoltageLightCtl       uint16  //光控电压    E01FH  0-14V
	SpecialCtl1           uint8   //特殊功率控制 E021H  b1: 1为开启特殊功 率控制功能，0为关闭特殊功率控制功能；b0: 1为开启每晚开启负载功能， 0为关闭每晚开启负载功能
	SpecialCtl2           uint8   //特殊功率控制        b2:零下禁充 1:开启零下禁止充电功能，0:关闭零下禁止充电功能 b0-b1:充电方式 00:直充充电方式， 01:PWM充电方式
}

type PowerTypeCapacity struct {
	BatterCapacity uint16 `json:"battercapacity"`
	BatterType     uint16 `json:"battertype"`
}

var FaultInfo = [32]string{"蓄电池过放", //0
	"蓄电池超压",     //1
	"欠压警告",      //2
	"负载短路",      //3
	"负载功率过大或过流", //4
	"控制器温度过高",   //5
	"电池高温保护（温度高于充电上限）禁止充电", //6
	"光伏输入功率过大",             //7
	"",                     //8
	"光伏输入端超压",              //9
	"",                     //10
	"太阳能板工作点超压",            //11
	"太阳能板反接",               //12
	"",                     //13
	"",                     //14
	"",                     //15
	"",                     //16
	"",                     //17
	"",                     //18
	"",                     //19
	"",                     //20
	"",                     //21
	"供电状态（市电供电）",           //22
	"未检测到电池（铅酸）",           //23
	"电池高温保护（温度高于放电上限）禁止放电", //24
	"电池低温保护（温度低于放电上限）禁止放电", //25
	"过冲保护，停止充电",            //26
	"电池低温保护（温度低于充电上限）禁止充电", //27
	"蓄电池反接",      //28
	"电容超压（保留）",   //29
	"感应探头损坏（路灯）", //30
	"负载开路（路灯）"}   //31

// var FaultInfo[0] := "ASD"
var powerreadbuff = make([]byte, 256)

func ReadPowerRunData(power SerialSet) (PowerRunData, error) {
	ret := PowerRunData{}

	port, err := serial.Open(SetSerialOption(power))
	if err != nil {
		return ret, err
	}
	defer port.Close()

	sendbuff := makeReadPowerRundata()
	n, errw := port.Write(sendbuff)
	if errw != nil {
		return ret, errw
	}

	if n < len(sendbuff) {
		return ret, errors.New("send cmd err")
	}

	time.Sleep(200 * time.Millisecond)

	len, errr := port.Read(powerreadbuff)
	if errr != nil {
		return ret, errr
	}
	if len < 75 {
		return ret, errors.New("read data err")
	}
	//fmt.Println(powerreadbuff[:len])
	data, eee := analysisPowerRunData(powerreadbuff[:len])
	if !eee {
		return ret, errors.New("analysys err")
	}
	//fmt.Println(data)
	return data, nil
}

func ReadPowerInfo(power SerialSet) (PowerInfo, error) {

	ret := PowerInfo{}

	port, err := serial.Open(SetSerialOption(power))
	if err != nil {
		//fmt.Println("serial open error")
		return ret, err
	}
	//fmt.Println(setSerialOption(power))
	//fmt.Println("open dev ok")

	defer port.Close()

	sendbuff := makeReadPowerinfo()
	//config.
	n, errw := port.Write(sendbuff)
	if errw != nil {
		return ret, errw
	}
	//fmt.Println(sendbuff)
	//fmt.Println("send ok!")

	if n < len(sendbuff) {
		return ret, errors.New("send cmd err")
	}
	//fmt.Println("send：", n)

	time.Sleep(200 * time.Millisecond)

	len, errr := port.Read(powerreadbuff)
	if errr != nil {
		return ret, errr
	}
	//fmt.Println(len)
	//fmt.Println(powerreadbuff[:len])
	// 01 03 34(字节数) xx ..34字节....xx crc1 crc2
	if len < 39 {
		return ret, errors.New("read data err")
	}
	info, eee := analysisPowerInfo(powerreadbuff[:len])
	if !eee {
		return ret, errors.New("analysys err")
	}
	//fmt.Println(info)
	return info, nil
}

func ReadPowerSet(power SerialSet) (PowerSet, error) {

	ret := PowerSet{}
	port, err := serial.Open(SetSerialOption(power))
	if err != nil {
		return ret, err
	}
	defer port.Close()

	sendbuff := makeReadPowerSet()
	n, errw := port.Write(sendbuff)
	if errw != nil {
		return ret, errw
	}
	if n < len(sendbuff) {
		return ret, errors.New("send cmd err")
	}

	time.Sleep(200 * time.Millisecond)

	len, errr := port.Read(powerreadbuff)
	if errr != nil {
		return ret, errr
	}
	//62+5=67
	if len < 67 {
		return ret, errors.New("read data err")
	}
	set, eee := analysisPowerSet(powerreadbuff[:len])
	if !eee {
		return ret, errors.New("analysys err")
	}

	return set, nil
}

func SetPowerBatterType(power SerialSet, t uint16) error {
	//ret := PowerSet{}
	if t > 4 {
		return errors.New("input para error")
	}
	port, err := serial.Open(SetSerialOption(power))
	if err != nil {
		return err
	}
	defer port.Close()

	sendbuff := makeSetPowerBatterType(t)
	n, errw := port.Write(sendbuff)
	if errw != nil {
		return errw
	}
	if n < len(sendbuff) {
		return errors.New("send cmd err")
	}
	time.Sleep(200 * time.Millisecond)

	len, errr := port.Read(powerreadbuff)
	if errr != nil {
		return errr
	}
	//62+5=67
	if len < 8 {
		return errors.New("recive data err")
	}
	if powerreadbuff[1] != 6 {
		return errors.New("recive data err")
	}

	return nil
}

func SetPowerBatterCapacity(power SerialSet, ah uint16) error {
	port, err := serial.Open(SetSerialOption(power))
	if err != nil {
		return err
	}
	defer port.Close()

	sendbuff := makeSetPowerBatterCapacity(ah)
	n, errw := port.Write(sendbuff)
	if errw != nil {
		return errw
	}
	if n < len(sendbuff) {
		return errors.New("send cmd err")
	}
	time.Sleep(200 * time.Millisecond)

	len, errr := port.Read(powerreadbuff)
	if errr != nil {
		return errr
	}
	//62+5=67
	if len < 8 {
		return errors.New("recive data err")
	}
	if powerreadbuff[1] != 6 {
		return errors.New("recive data err")
	}

	return nil

}

func makeReadPowerinfo() []byte {
	//读取
	//|01 |02    |03-04  |05-06   |07-08|
	//|地址|功能码|起始地址|长度（字）|CRC16|
	sendbuff := make([]byte, 8)
	sendbuff[0] = 1 //设备地址
	sendbuff[1] = 3 //读取指令
	sendbuff[2] = 0
	sendbuff[3] = 0x0A
	sendbuff[4] = 0
	sendbuff[5] = 0x11 //17字=34字节
	//fmt.Println(sendbuff[:6])
	crc16 := makeCRC16(sendbuff[:6])
	sendbuff[6] = byte(crc16 & 0xff)
	sendbuff[7] = byte(crc16 >> 8)

	//fmt.Println(checkCRC16(sendbuff))
	return sendbuff

}

func makeReadPowerRundata() []byte {
	//读取
	//|01 |02    |03-04  |05-06   |07-08|
	//|地址|功能码|起始地址|长度（字）|CRC16|
	sendbuff := make([]byte, 8)
	sendbuff[0] = 1 //设备地址
	sendbuff[1] = 3 //读取指令
	sendbuff[2] = 0x01
	sendbuff[3] = 0x00
	sendbuff[4] = 0
	sendbuff[5] = 0x23 //35字=70字节
	crc16 := makeCRC16(sendbuff[:6])
	sendbuff[6] = byte(crc16 & 0xff)
	sendbuff[7] = byte(crc16 >> 8)
	return sendbuff
}

func makeReadPowerSet() []byte {
	sendbuff := make([]byte, 8)
	sendbuff[0] = 1 //设备地址
	sendbuff[1] = 3 //读取指令
	sendbuff[2] = 0xE0
	sendbuff[3] = 0x02
	sendbuff[4] = 0
	sendbuff[5] = 0x1F //31字=62字节
	crc16 := makeCRC16(sendbuff[:6])
	sendbuff[6] = byte(crc16 & 0xff)
	sendbuff[7] = byte(crc16 >> 8)
	return sendbuff

}

func makeSetPowerBatterType(t uint16) []byte {
	sendbuff := make([]byte, 8)
	sendbuff[0] = 1 //设备地址
	sendbuff[1] = 6 //写入指令
	sendbuff[2] = 0xE0
	sendbuff[3] = 0x04
	sendbuff[4] = uint8(t >> 8) //
	sendbuff[5] = uint8(t & 0xFF)
	crc16 := makeCRC16(sendbuff[:6])
	sendbuff[6] = byte(crc16 & 0xff)
	sendbuff[7] = byte(crc16 >> 8)
	return sendbuff
}

func makeSetPowerBatterCapacity(ah uint16) []byte {
	sendbuff := make([]byte, 8)
	sendbuff[0] = 1 //设备地址
	sendbuff[1] = 6 //写入指令
	sendbuff[2] = 0xE0
	sendbuff[3] = 0x02
	sendbuff[4] = uint8(ah >> 8) //
	sendbuff[5] = uint8(ah & 0xFF)
	crc16 := makeCRC16(sendbuff[:6])
	sendbuff[6] = byte(crc16 & 0xff)
	sendbuff[7] = byte(crc16 >> 8)
	return sendbuff
}

func analysisPowerInfo(buff []byte) (PowerInfo, bool) {
	ret := PowerInfo{}

	if !checkCRC16(buff) {
		return ret, false
	}
	if len(buff) != 34+5 {
		return ret, false
	}
	datelen := int(buff[2])
	if datelen != 34 {
		return ret, false
	}

	//最大支持电压、额定充电电流、额定放电电流、产品类型、型号、软件版本、硬件版本、序列号
	//0A                    0B                         0C -13         |14-15        |16-17      |18-19   |20   |
	//|01        |02        |03        |04             | 05-20（8字）  |21-24        |25-28      |29-32   |33-34|
	//最大支持电压 |额定充电电流|额定放电电流|产品类型         |型号          |软件版本      |硬件版本     |序列号   |地址  |
	//10进制+V    |10进制+A   |10进制+A  |00控制器、02逆变器|（ascii码     |后三字节每字节代表一个版本子号|转成hex串|低字节 |
	data := buff[3 : len(buff)-2]
	ret.MaxSuppoptVoltage = fmt.Sprint(data[0], "V")
	ret.RateChargeCurrent = fmt.Sprint(data[1], "A")
	ret.RateDischargeCurrent = fmt.Sprint(data[2], "A")
	if data[3] == 0 {
		ret.ProductType = "控制器"
	} else {
		ret.ProductType = "逆变器"
	}
	ret.ProductModel = string(data[4:20])
	ret.SoftwareVer = fmt.Sprint(data[21], ".", data[22], ".", data[23])
	ret.HardwareVer = fmt.Sprint(data[25], ".", data[26], ".", data[27])
	ret.SN = hex.EncodeToString(data[28:32])

	return ret, true
}

func analysisPowerRunData(buff []byte) (PowerRunData, bool) {
	ret := PowerRunData{}

	if !checkCRC16(buff) {
		return ret, false
	}
	if len(buff) != 70+5 {
		return ret, false
	}

	datelen := int(buff[2])
	if datelen != 70 {
		return ret, false
	}

	//0100H------------------------------70字节-----------------------------------------------------------0120H
	//0100-------------------------------14字节--------------------------------------------------------------0106
	//|01-02             |03-04         |05-06      |07                      |08     |09-10        |11-12           |13-14   |
	//|电池电量02（0-100）%|蓄电池电压V 0.1|充电电流A0.01|控制器温度（b7符号，b0-6值）|电池温度|负载直流电压V0.1|负载直流电流A0.01|负载功率W|
	data := buff[3 : 3+14]

	ret.BatterySOC = uint(data[1])
	ret.BatteryVoltage = MyFloat(computeint16(data[2:4])) * 0.1
	ret.ChargeCurrent = MyFloat(computeint16(data[4:6])) * 0.01

	ret.TempControl = int(data[6] & 0x7F)
	if data[6]&0x80 == 0x80 {
		ret.TempControl = -ret.TempControl
	}

	ret.TempBattery = int(data[7] & 0x7F)
	if data[7]&0x80 == 0x80 {
		ret.TempBattery = -ret.TempBattery
	}

	ret.LoadVoltage = MyFloat(computeint16(data[8:10])) * 0.1
	ret.LoadCurrent = MyFloat(computeint16(data[10:12])) * 0.01
	ret.LoadPower = computeint16(data[12:14])

	//太阳能板电压，太阳能板电流，充电功率
	//0107-------------------------------太阳能板信息（6 字节）--------------------------0109H
	//|01-02     |03-04     |05-06  |
	//|太阳能板电压|太阳能板电流|充电功率|
	solardata := buff[3+14 : 3+14+6]
	ret.SolarVoltage = MyFloat(computeint16(solardata[0:2])) * 0.1
	ret.SolarCurrent = MyFloat(computeint16(solardata[2:4])) * 0.01
	ret.SolarPower = computeint16(solardata[4:6])
	//开/关负载命令,蓄电池当天最低电压,蓄电池当天最高电压,当天充电最大电流,当天放电最大电流,当天充电最大功率,当天放电最大功率,当天充电安时数,当天放电安时数,当天发电量,当天用电量
	//010AH-----------------------------蓄电池信息（22字节）-----------------------------0114H
	//|01-02                       |03-04            |05-06  （电压*0.1，电流*0.01，功率原值）
	//|开/关负载命令0000关闭，0001打开|蓄电池当天最低电压0.1|

	batterydata := buff[3+14+6 : 3+14+6+22]
	if batterydata[1] == 1 {
		ret.LoadSwitch = true
	} else {
		ret.LoadSwitch = false
	}

	ret.BatterVoltageMin = MyFloat(computeint16(batterydata[2:4])) * 0.1
	ret.BatterVoltageMax = MyFloat(computeint16(batterydata[4:6])) * 0.1
	ret.ChargeCurrentMax = MyFloat(computeint16(batterydata[6:8])) * 0.01
	ret.DischargeCurrentMax = MyFloat(computeint16(batterydata[8:10])) * 0.01
	ret.ChargePowerMax = computeint16(batterydata[10:12])
	ret.DischargePowerMax = computeint16(batterydata[12:14])
	ret.ChargeAH = computeint16(batterydata[14:16])
	ret.DischargeAH = computeint16(batterydata[16:18])
	ret.GenPower = computeint16(batterydata[18:20])
	ret.UsePower = computeint16(batterydata[20:22])

	//0115H-----------------------------历史数据信息（22字节）--------------------------011FH
	//总运行天数，蓄电池总过放次数，蓄电池总充满次数，蓄电池总充电安时数（4b），蓄电池总放电安时数（4b)，累计发电量（4b），累计用电量（4b）
	histortdata := buff[3+14+6+22 : 3+14+6+22+22]
	ret.RunDaysCount = computeint16(histortdata[0:2])
	ret.OverDischargeCount = computeint16(histortdata[2:4])
	ret.FillingsCount = computeint16(histortdata[4:6])
	ret.ChargeAHCount = uint32(computeint32(histortdata[6:10]))
	ret.DischargeAHCount = uint32(computeint32(histortdata[10:14]))
	ret.GenPowerCount = uint32(computeint32(histortdata[14:18]))
	ret.UsePowerCount = uint32(computeint32(histortdata[18:22]))

	//0120
	//负载状态，负载亮度，充电状态
	//负载状态 高字节b7 0 负载关 1 负载开
	//负载亮度 高字节b0-6
	//充电状态：00H:未开启充电 01H:启动充电模式 02 mppt充电模式  03H:均衡充电模式 04H:提升充电模式 05H:浮充充电模式 06H:限流(超功率)
	//----------------------------------------负载信息（2字节）------------------------
	loaddata := buff[3+14+6+22+22 : 3+14+6+22+22+2]
	if loaddata[0]&0x80 == 0x80 {
		ret.LoadStatus = true
	} else {
		ret.LoadStatus = false
	}

	ret.LoadRright = int(loaddata[0] & 0x7f)
	ret.ChargeStatus = int(loaddata[1])

	//0121-0122
	//故障码总计32位代表21种错误，其他保留
	faultdata := buff[3+14+22+22+22+2 : 3+14+22+22+22+2+4]
	ret.FaultCode = computeint32(faultdata[0:4])
	ret.FaultStrings = makeFaultStrings(ret.FaultCode)

	return ret, true
}

func makeFaultStrings(code uint32) []string {
	var ret []string
	for i := 0; i < 32; i++ {
		if code>>i&0x01 == 0x01 {
			ret = append(ret, FaultInfo[i])
		}
	}
	//len(ret)
	return ret
}

func analysisPowerSet(buff []byte) (PowerSet, bool) {
	ret := PowerSet{}

	if !baseCheck(buff) {
		return ret, false
	}

	data := buff[3:]
	ret.BatterCapacity = computeint16(data[0:2])
	ret.VoltageSet = data[2]
	ret.VoltageIden = data[3]
	ret.BatterType = computeint16(data[4:6])
	ret.VoltageOver = MyFloat(computeint16(data[6:8])) * 0.1
	ret.VoltageChargeLim = MyFloat(computeint16(data[8:10])) * 0.1
	ret.VoltageChargeEqu = MyFloat(computeint16(data[10:12])) * 0.1
	ret.VoltageChargeUpg = MyFloat(computeint16(data[12:14])) * 0.1
	ret.VoltageChargeFlo = MyFloat(computeint16(data[14:16])) * 0.1
	ret.VoltageChargeUpgR = MyFloat(computeint16(data[16:18])) * 0.1
	ret.VoltageDischargeOverR = MyFloat(computeint16(data[18:20])) * 0.1
	ret.VoltageUnderWarn = MyFloat(computeint16(data[20:22])) * 0.1
	ret.VoltageDischargeOver = MyFloat(computeint16(data[22:24])) * 0.1
	ret.VoltageDisChargeLim = MyFloat(computeint16(data[24:26])) * 0.1
	ret.SocChargeLim = data[26]
	ret.SocDischargeLim = data[27]
	ret.DelaySDischargeOver = computeint16(data[28:30])
	ret.TimeMChargeEqu = computeint16(data[30:32])
	ret.TimeMChargeUpg = computeint16(data[32:34])
	ret.IntervalDChargeUpg = computeint16(data[34:36])
	ret.TempComp = computeint16(data[36:38])

	ret.LoadMode = computeint16(data[52:54])
	ret.DelayMLightCtl = computeint16(data[54:56])
	ret.VoltageLightCtl = computeint16(data[56:68])
	ret.SpecialCtl1 = data[60]
	ret.SpecialCtl2 = data[61]

	return ret, true
}

func baseCheck(buff []byte) bool {
	if !checkCRC16(buff) {
		return false
	}
	//01 02 03 ----------------------------//
	//01 02 len-------------------CRC1 CRC2//
	return int(buff[2]) == len(buff)-5
}
