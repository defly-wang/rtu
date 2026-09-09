package config

import (
	"bufio"
	"os"
	"strings"
)

type NetworkConfig struct {
	Ip      string `json:"ip"`
	Mask    string `json:"mask"`
	Gateway string `json:"gateway"`
}

func GetNetworkConfig() (NetworkConfig, error) {
	var info NetworkConfig

	file, err := os.Open("/mnt/yaffs2/net.conf")
	if err != nil {
		return info, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "IPADDR=") {
			info.Ip = strings.TrimSpace(strings.TrimPrefix(line, "IPADDR="))
		} else if strings.HasPrefix(line, "NETMASK=") {
			info.Mask = strings.TrimSpace(strings.TrimPrefix(line, "NETMASK="))
		} else if strings.HasPrefix(line, "GATEWAY=") {
			info.Gateway = strings.TrimSpace(strings.TrimPrefix(line, "GATEWAY="))
		}
	}

	return info, nil
}

func SaveNetworkConfig(cfg NetworkConfig) bool {
	file, err := os.Create("/mnt/yaffs2/net.conf")
	if err != nil {
		return false
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
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

	err = writer.Flush()
	if err != nil {
		return false
	}

	return true
}
