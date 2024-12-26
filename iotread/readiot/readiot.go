package readiot

import (
	"config"
	"encoding/hex"
	"errors"
	"fmt"
	"iotread/common"

	"time"

	"github.com/jacobsa/go-serial/serial"
)

// 基准值
var baseData map[int]config.IotInfo = make(map[int]config.IotInfo)

var recivebuff = make([]byte, 256)

func ReadIot(iot int) (config.IotInfo, error) {
	//读取等待延时(毫秒),读取长度,数据长度:1\2\4
	//serial.GetPorts()

	options := serial.OpenOptions{
		PortName:              common.GetIotPortName(iot),                         //config.CFG.Iots[iot].Com
		BaudRate:              common.CFG.Iots[iot].Baudrate,                      //config.GetIotBaudRate(iot),
		DataBits:              common.CFG.Iots[iot].Databits,                      //config.GetIotDataBits(iot),
		StopBits:              common.CFG.Iots[iot].Stopbits,                      //config.GetIotStopBits(iot),
		ParityMode:            serial.ParityMode(common.CFG.Iots[iot].Paritymode), //config.GetIotParityMode(iot)),
		InterCharacterTimeout: 200,                                                //此参数必须设置，否则会长时间等待
	}

	iotinfo := common.InitIotInfo(iot)
	port, err := serial.Open(options)
	if err != nil {
		fmt.Println("serial open error")
		return iotinfo, err
	}

	defer port.Close()

	//port.

	//发送指令,长度 例如 040300000001845f
	//fmt.Println(GetIotReadbuff(iot))
	//此指令未测试是否有必要
	//port.Read(nil)
	n, err := port.Write([]byte(common.GetIotReadbuff(iot))) //config.GetIotReadbuff(iot))
	if err != nil {
		//fmt.Println("send error")
		return iotinfo, err
	}
	//fmt.Println(config.CFG.Iots[iot].Read_len)

	if n != common.CFG.Iots[iot].Read_len {
		//fmt.Println(n)
		return iotinfo, errors.New("send error2")
	}

	time.Sleep(time.Duration(common.CFG.Iots[iot].Read_delay) * time.Millisecond)

	rec_len, err := port.Read(recivebuff)
	if err != nil {
		return iotinfo, err
	}

	if rec_len < common.CFG.Iots[iot].Recive_len {
		return iotinfo, errors.New("read error2")
	}

	if !check(iot, recivebuff[:rec_len]) {
		return iotinfo, errors.New("CRC or Xor error")
	}

	iotinfo.Rawdata = hex.EncodeToString(recivebuff[:rec_len])

	//var result uint64
	//计算需要重写（各种类型）
	compute(iot, recivebuff, &iotinfo)

	return iotinfo, nil
}

func ReadIotJson(iot int, times int) (config.IotInfo, error) {

	for i := 0; i < times; i++ {
		data, err := ReadIot(iot)
		if err == nil {
			return data, err
		}
	}
	return config.IotInfo{}, errors.New(fmt.Sprint("read error ", times, "times"))
}

// 检查校验
func check(iot int, buff []byte) bool {
	switch common.CFG.Iots[iot].Crc {
	case "CRC16":
		if !checkCRC16(buff) {
			return false
		}
	case "XOR":
		if !checkXOR(buff) {
			return false
		}
	}
	return true
}

// CRC16校验 返回 低字节---高字节
func checkCRC16(buff []byte) bool {
	//crc32.c
	//var crc16 uint16
	//uint16_t j = 0;
	crc16 := uint16(0xFFFF)

	for i := 0; i < len(buff); i++ {
		crc16 ^= uint16(buff[i])
		for j := 0; j < 8; j++ {
			if crc16&0x01 != 0 {
				crc16 >>= 1
				crc16 ^= 0xA001 //0xA001为0x8005按bit位颠倒后的生成项
			} else {
				crc16 >>= 1
			}

		}
	}
	return crc16 == 0x00
}

// checkXOR
func checkXOR(buff []byte) bool {
	crc := uint8(0x00)
	for i := 0; i < len(buff); i++ {
		crc ^= buff[i]
	}
	return crc == 0x00
}

func compute(iot int, buff []byte, iotinfo *config.IotInfo) {
	iotinfo.Value = 0

	for i := 0; i < common.CFG.Iots[iot].Revive_data_len; i++ {
		//revive_data1 ---revive_data4
		//各个字节
		//fmt.Println("i", i, j, revivebuff[j])
		iotinfo.Value = (int64)(iotinfo.Value*256) + (int64)(buff[common.CFG.Iots[iot].Revive_data[i]])
	}
	if common.CFG.Iots[iot].Offset { //} IsIotCalOffset(iot) {
		oldiotinfo := FindMapIot(iot, *iotinfo)
		iotinfo.Offset = iotinfo.Value - oldiotinfo.Value
	}

	iotinfo.Fvalue = config.MyFloat(iotinfo.Value) * config.MyFloat(iotinfo.Ratio)
	iotinfo.Foffset = config.MyFloat(iotinfo.Offset) * config.MyFloat(iotinfo.Ratio)

	//save.Save(iotinfo)

}

// 基准值
func FindMapIot(iot int, newiotinfo config.IotInfo) config.IotInfo {

	//查找
	iotinfo, ok := baseData[iot]
	//更新
	baseData[iot] = newiotinfo
	if ok {
		return iotinfo
	} else {
		//第一次，返回原值
		return newiotinfo
	}
}
