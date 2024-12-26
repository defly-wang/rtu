package save

import (
	"config"
	"encoding/json"
	"fmt"
	"iotread/common"
	"os"
	"time"
)

// var datadir string
const (
	// format = "2006-01-02 15:04:05"
	tfmt = "15:04:05"

// dfmt = "2006-01-02"
)

func GetDataDir(t time.Time) string {
	return fmt.Sprintf("%s/%4d/%02d/", common.CFG.Common.Datadir, t.Year(), int(t.Month()))
}

func GetSaveFileName(t time.Time, iot int) string {
	//sid := hex.EncodeToString([]byte(id))
	return fmt.Sprintf("%s/IOT-%02d-%02d.txt", GetDataDir(t), iot, t.Day())
}

func GetUnSendFileName() string {

	return fmt.Sprintf("%s/unsended.txt", common.CFG.Common.Datadir)
}

func isDirExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func Save(data config.IotInfo) {

	t := time.Time(data.Time)

	path := GetDataDir(t)
	if !isDirExists(path) {
		//fmt.Println("目录不存在！")
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return
			//fmt.Println(err.Error())
		}

	}

	file, err := os.OpenFile(GetSaveFileName(t, data.Iot), os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		//fmt.Printf("Error:%s", err.Error())
		return
	}
	defer file.Close()

	datajson, jerr := json.Marshal(data)
	if jerr != nil {
		//fmt.Println("json:", jerr.Error())
		return
	}

	//"%02d-%02d-%02d,%s\n", t.Hour(), t.Minute(), t.Second()
	file.WriteString(fmt.Sprintf("%s,%s\n", t.Format(tfmt), string(datajson)))

	//file.Write([]byte(w))
}

/*
	func Load(year int, month int, iot int) []config.IotConfig {
		result := []config.IotConfig{}

		t := time.Date(year, time.Month(0), 0, 0, 0, 0, 0, time.UTC)
		filename := GetSaveFileName(t, iot)

		file, _ := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
		reader := bufio.NewReader(file)

		var iotinfo config.IotConfig
		for {
			bytes, _, err := reader.ReadLine()
			if err == io.EOF {
				break
			}

			errj := json.Unmarshal(bytes[10:], &iotinfo)
			if errj != nil {
				fmt.Println(err.Error())
			}
			result = append(result, iotinfo)

		}
		return result
	}
*/

/*
	func LoadOne(t time.Time, iot int) config.IotInfo {
		filename := GetSaveFileName(t, iot)
		file, _ := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
		reader := bufio.NewReader(file)
		var diff time.Duration = 1000000
		var iotinfo config.IotInfo

		for {
			bytes, _, err := reader.ReadLine()
			if err == io.EOF {
				break
			}

			tt, error := time.Parse(format, string(bytes)[:10])
			if error != nil {
				return iotinfo
			}
			if t.Sub(tt) < diff {
				diff = t.Sub(tt)
				json.Unmarshal(bytes[10:], &iotinfo)
			}

		}

		return iotinfo

}
*/
/*
func Cont(year int, month int, iot int) int {

	t := time.Date(year, time.Month(0), 0, 0, 0, 0, 0, time.UTC)
	filename := GetSaveFileName(t, iot)

	file, _ := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	fmt.Println(file)
	return 0
}

func CountAll(iot int) map[string]int {
	//from :=
	//三年
	var result = make(map[string]int)
	from := time.Now().AddDate(-3, 0, 0)
	for {
		cont := Cont(from.Year(), int(from.Month()), iot)
		result[from.Format("2006-01-02")] = cont
		from := from.AddDate(0, 1, 0)
		if from.After(time.Now()) {
			break
		}
		//t1 := time.Parse(format,fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",f))

	}
	return result

}
*/
