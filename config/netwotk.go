package config

import (
	"bufio"
	"os"
	"strings"
)

type NetworkConfig struct {
	Method  string
	Ip      string `json:"ip"`
	Mask    string `json:"mask"`
	Gateway string `json:"gateway"`
	Mac     string
}

func GetNetworkConfig() (NetworkConfig, error) {
	var info NetworkConfig

	file, err := os.Open("/mnt/yaffs2/net.conf")
	if err != nil {
		return info, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	info.Method = scanLine(scanner, "METHOD=")
	info.Ip = scanLine(scanner, "IPADDR=")
	info.Mask = scanLine(scanner, "NETMASK=")
	info.Gateway = scanLine(scanner, "GATEWAY=")
	info.Mac = scanLine(scanner, "MAC=")

	return info, nil
}

func scanLine(scanner *bufio.Scanner, key string) string {
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, key) {
			return strings.TrimPrefix(line, key)
		}
	}
	return ""
}

func SaveNetworkConfig(cfg NetworkConfig) bool {

	oldCfg, err := GetNetworkConfig()
	if err != nil {
		return false
	}

	file, err := os.Create("/mnt/yaffs2/net.conf")
	if err != nil {
		return false
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	_, err = writer.WriteString("METHOD=" + oldCfg.Method + "\n")
	if err != nil {
		return false
	}

	_, err = writer.WriteString("IPADDR=" + cfg.Ip + "\n")
	if err != nil {
		return false
	}
	_, err = writer.WriteString("NETMASK=" + cfg.Mask + "\n")
	if err != nil {
		return false
	}
	_, err = writer.WriteString("GATEWAY=" + cfg.Gateway + "\n")
	if err != nil {
		return false
	}

	_, err = writer.WriteString("MAC=" + oldCfg.Mac + "\n")
	if err != nil {
		return false
	}

	err = writer.Flush()
	if err != nil {
		return false
	}

	return true
}
