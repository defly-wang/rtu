//go:build !windows

package config

import (
	"bytes"
	"errors"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Sysinfo struct {
	Version       string `json:"version"`
	Iotread_ver   string `json:"read_version"`
	Rtuclient_ver string `json:"rtu_version"`
	Webcgi_ver    string `json:"cgi_version"`
	Uptime        int64  `json:"uptime"`
	Memtotal      uint64 `json:"memtotal"`
	Memfree       uint64 `json:"memfree"`
	Disktotal     uint64 `json:"disktotal"`
	DiskFree      uint64 `json:"diskfree"`
	Mac           string `json:"mac"`
	Ip            string `json:"ip"`
	Gprs_ip       string `json:"gprs_ip"`
	Gpis_imsi     string `json:"gprs_imsi"`
	Gpis_ccid     string `json:"gprs_ccid"`
	Gprs_csq      string `json:"gprs_csq"`
}

func GetSysInfo() Sysinfo {
	var info Sysinfo

	interfaces, err := net.Interfaces()
	if err != nil {
		return info
	}

	for _, inter := range interfaces {
		if inter.Name == "eth0" {
			info.Mac = inter.HardwareAddr.String()
			add, _ := inter.Addrs()
			for _, ad := range add {
				ipnet, _ := ad.(*net.IPNet)
				if ipnet.IP.To4() != nil {
					info.Ip = ad.String()
				}
			}
		}
		if inter.Name == "usb0" {
			add, _ := inter.Addrs()
			for _, ad := range add {
				ipnet, _ := ad.(*net.IPNet)
				if ipnet.IP.To4() != nil {
					info.Gprs_ip = ad.String()
				}
			}
		}

	}

	info.Version, _ = getLinuxVersion()
	info.Iotread_ver = GetIotVersion()
	info.Webcgi_ver = GetCgiVersion()
	info.Rtuclient_ver = GetRtuClientVersion()
	//info.Rtuclient_ver =

	info.Gpis_imsi = Get4Ginfo(CMD_GET_IMSI)
	time.Sleep(100 * time.Millisecond)
	info.Gpis_ccid = Get4Ginfo(CMD_GET_CCID)
	time.Sleep(100 * time.Millisecond)
	info.Gprs_csq = Get4Ginfo(CMD_GET_CSQ)

	var Mem syscall.Sysinfo_t
	var Disk syscall.Statfs_t

	//var Mem runtime.MemStats
	//runtime.ReadMemStats(&Mem)
	//info.Memtotal = Mem.TotalAlloc
	//info.Memfree = Mem.Frees
	//runtime.Version()

	//m.Frees

	if syscall.Sysinfo(&Mem) != nil {
		return info
	}

	info.Memtotal = uint64(Mem.Totalram)
	info.Memfree = uint64(Mem.Freeram)
	info.Uptime = int64(Mem.Uptime)

	//info.Uptime, _ = getUptime()

	derr := syscall.Statfs(SYSTEM_CONFIG_DATAPATH, &Disk)
	if derr != nil {
		//fmt.Println(derr)
		return info
	}

	info.Disktotal = Disk.Blocks * uint64(Disk.Bsize)
	info.DiskFree = Disk.Bfree * uint64(Disk.Bsize)

	return info
}

func getLinuxVersion() (string, error) {
	var out bytes.Buffer
	cmd := exec.Command("uname", "-r", "-s", "-m")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	out.Truncate(out.Len() - 1)
	return out.String(), nil
}

func GetUptime() (int64, error) {

	connect, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0, err
	}
	strtime := strings.Split(string(connect), " ")
	if strtime == nil {
		return 0, errors.New("read error")
	}
	fuptime, errc := strconv.ParseFloat(strtime[0], 64)
	if errc != nil {
		return 0, errc
	}
	return int64(fuptime), nil
}
