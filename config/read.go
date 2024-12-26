package config

import (
	"github.com/jacobsa/go-serial/serial"
)

type SerialSet struct {
	Com        string `json:"com"`
	Baudrate   uint   `json:"baudrate"`
	Databits   uint   `json:"databits"`
	Stopbits   uint   `json:"stopbits"`
	Paritymode uint   `json:"paritymode"`
}

func SetSerialOption(ser SerialSet) serial.OpenOptions {
	return serial.OpenOptions{
		PortName:              ser.Com,
		BaudRate:              ser.Baudrate,
		DataBits:              ser.Databits,
		StopBits:              ser.Stopbits,
		ParityMode:            serial.ParityMode(ser.Paritymode),
		InterCharacterTimeout: 200, //此参数必须设置，否则会长时间等待
	}
}

func compute(iotcfg IotConfig, buff []byte, iotinfo *IotInfo) {
	iotinfo.Value = 0

	for i := 0; i < iotcfg.Revive_data_len; i++ {
		//revive_data1 ---revive_data4
		//各个字节
		//fmt.Println("i", i, j, revivebuff[j])
		iotinfo.Value = (int64)(iotinfo.Value*256) + (int64)(buff[iotcfg.Revive_data[i]])
	}

	//if iotcfg.Offset { //} IsIotCalOffset(iot) {
	//oldi := 0 //findMapIot(iot, *iotinfo)
	//iotinfo.Offset = iotinfo.Value - oldiotinfo.Value
	//}

	iotinfo.Fvalue = MyFloat(iotinfo.Value) * MyFloat(iotinfo.Ratio)
	//iotinfo.Foffset = config.MyFloat(iotinfo.Offset) * config.MyFloat(iotinfo.Ratio)

}

func computeint16(buff []byte) uint16 {
	if len(buff) < 2 {
		return 0
	}
	var ret uint16 = 0
	ret = uint16(buff[0])<<8 + uint16(buff[1])
	return ret
}
func computeint32(buff []byte) uint32 {
	if len(buff) < 4 {
		return 0
	}
	var ret uint32 = 0
	ret = uint32(buff[0])<<24 + uint32(buff[1])<<16 + uint32(buff[2])<<8 + uint32(buff[3])
	return ret
}
