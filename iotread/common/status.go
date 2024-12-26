package common

import (
	"config"
	"encoding/json"
	"sync"
	"time"
)

var runstatus = config.Status{
	Starttime:   config.MyTime(time.Now()),
	Count:       0,
	Readcount:   0,
	Mqttcount:   0,
	Webapicount: 0,
	Runing:      false,
	Version:     VERSION,
}

var lock sync.RWMutex

func SetStatusRuning(run bool) {
	lock.Lock()
	runstatus.Runing = run
	lock.Unlock()
}

func StatusCount() {
	lock.Lock()
	runstatus.Count++
	lock.Unlock()
}

func StatusReadSuccess() {
	lock.Lock()
	runstatus.Readcount++
	lock.Unlock()
}

func StatusMqttSuccess() {
	lock.Lock()
	runstatus.Mqttcount++
	lock.Unlock()
}

func StatusWebapiSuccess() {
	lock.Lock()
	runstatus.Webapicount++
	lock.Unlock()
}

func GetStatus() config.Status {
	lock.RLock()
	ret := runstatus
	lock.RUnlock()
	return ret
}

func GetStatusJson() []byte {
	ret, err := json.Marshal(GetStatus())
	if err != nil {
		return nil
	}
	return ret
}
