package config

var (
	CRC_CRC16 = "CRC16"
	CRC_XOR   = "XOR"
)

func CheckCRC(buff []byte, t string) bool {
	switch t {
	case CRC_CRC16:
		return checkCRC16(buff)
	case CRC_XOR:
		return checkXOR(buff)
	default:
		return false
	}
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

func makeCRC16(buff []byte) uint16 {

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
	return crc16

}

// checkXOR
func checkXOR(buff []byte) bool {
	crc := uint8(0x00)
	for i := 0; i < len(buff); i++ {
		crc ^= buff[i]
	}
	return crc == 0x00
}
