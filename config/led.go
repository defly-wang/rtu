package config

import (
	"os"
	"time"
)

const (
	LedFile = "/dev/led"
)

func setLed(id int, on bool) bool {
	filename := LedFile
	if id == 2 {
		filename = filename + "2"
	}
	file, err := os.OpenFile(filename, os.O_RDWR, 0)
	if err != nil {
		return false
	}
	defer file.Close()

	buf := make([]byte, 1)
	if on {
		buf[0] = '1'
	} else {
		buf[0] = '0'
	}
	_, errw := file.Write(buf)

	return errw == nil
}

func getLed(id int) bool {
	filename := LedFile
	if id == 2 {
		filename = filename + "2"
	}
	file, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return false
	}
	defer file.Close()

	buf := make([]byte, 10)
	n, errf := file.Read(buf)

	if n < 0 || errf != nil {
		return false
	}
	return buf[0] == '1'
	//return false
}

func LedOn() bool {
	return setLed(1, true)
}

func LedOff() bool {
	return setLed(1, false)
}

func Led2On() bool {
	return setLed(2, true)
}

func Led2Off() bool {
	return setLed(2, false)
}

func ledBlink(id int, t uint) {
	on := getLed(id)
	defer setLed(id, on)

	now := time.Now().Add(time.Millisecond * time.Duration(t))
	for {
		if time.Now().After(now) {
			return
		}
		setLed(id, true)
		time.Sleep(100 * time.Millisecond)
		setLed(id, false)
		time.Sleep(100 * time.Millisecond)
	}
}

func LedBlink(t uint) {
	ledBlink(1, t)
}

func LedBlink2(t uint) {
	ledBlink(2, t)
}
