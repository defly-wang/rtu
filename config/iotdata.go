package config

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	format = "2006-01-02 15:04:05"
	dfmt   = "2006-01-02"
	tfmt   = "15:04:05"
)

type MyTime time.Time
type MyFloat float64

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type IotInfo struct {
	Iot     int     `json:"-"`
	Id      string  `json:"id"`
	Type    string  `json:"type"`
	Value   int64   `json:"-"`
	Offset  int64   `json:"-"`
	Ratio   float64 `json:"-"`
	Fvalue  MyFloat `json:"value"`
	Foffset MyFloat `json:"offset"`
	Rawdata string  `json:"rawdata"`
	Time    MyTime  `json:"time"`
	// Ok      bool
}

//var baseData map[int]IotInfo = make(map[int]IotInfo)

// Json 序列化
func (t MyTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(format)+2)
	b = append(b, '"')
	b = (time.Time(t)).AppendFormat(b, format)
	b = append(b, '"')
	return b, nil
}

// Json 反序列化
func (t *MyTime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+format+`"`, string(data), time.Local)
	*t = MyTime(now)
	return
}

// MarshalJSON 实现了json.Marshaler接口，用于自定义序列化行为
func (f MyFloat) MarshalJSON() ([]byte, error) {
	str := strconv.FormatFloat(float64(f), 'f', 2, 64) // 2表示小数点后保留2位
	return []byte(str), nil
}

/*
	func SetRawData(buff []byte, len int) []byte {
		raw := make([]byte, len)
		copy(raw, buff)
		return raw
	}
*/
func IotDataToJson(data IotInfo) []byte {
	jsondata, err := json.Marshal(data)
	if err != nil {
		return nil
	}
	return jsondata
}

func IotDataToJsonResult(data IotInfo, err error) []byte {

	var result Result
	if err != nil {
		result.Code = 203
		result.Msg = fmt.Sprintf("读取错误：%s", err.Error())
	} else {
		result.Code = 200
		result.Msg = "OK"
		result.Data = data
	}

	jsondata, err := json.Marshal(result)
	if err != nil {
		return nil
	}
	return jsondata
}

func StrToTime(str string) (time.Time, error) {
	var ret time.Time
	switch len(str) {
	case 10: //日期
		return time.Parse(dfmt, str)

	case 8: //时间
		return time.Parse(tfmt, str)

	case 19: //日期+时间
		return time.Parse(format, str)
	default:
		return ret, errors.New("time format error")

	}
}

func LoadHistoryData(cfg Config, t time.Time, iot int) ([]IotInfo, error) {
	result := []IotInfo{}
	//time 忽略 时分秒
	//year int, month int, day int,

	//t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	filename := GetSaveFileName(cfg, t, iot)

	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	var iotinfo IotInfo
	for {
		bytes, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		//fmt.Println(string(bytes[9:]))
		errj := json.Unmarshal(bytes[9:], &iotinfo)
		if errj != nil {
			fmt.Println(err)
			continue
		}
		result = append(result, iotinfo)

	}
	return result, nil
}

func LoadHistoryDataOne(cfg Config, t time.Time, iot int) (IotInfo, error) {
	//time 年月日 时分秒
	var iotinfo IotInfo
	var errret error = nil
	diff := 24 * time.Hour

	//fmt.Println("diff", diff)

	filename := GetSaveFileName(cfg, t, iot)

	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return iotinfo, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		bytes, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		/*
			if err == io.EOF {
				break
			}
		*/

		datatime, errt := time.Parse(format, t.Format(dfmt)+" "+string(bytes[0:8]))
		if errt != nil {
			//fmt.Println(errt)
			errret = errt
			break
		}

		var sub time.Duration
		if t.After(datatime) {
			sub = t.Sub(datatime)
		} else {
			sub = datatime.Sub(t)
		}
		//fmt.Println(sub)

		if sub < diff {
			diff = sub
			errj := json.Unmarshal(bytes[9:], &iotinfo)
			if errj != nil {
				errret = errj
				fmt.Println(err)
				continue
			}
		}
	}
	return iotinfo, errret
}
