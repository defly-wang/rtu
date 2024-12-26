package config

import (
	"os/exec"
	"strings"
)

const (
	CgiFile = "rtucgi.cgi"
	IotFile = "iotread"
	RtuFile = "rtuclient"
)

func GetIotVersion() string {
	return GetRtuRunStatus().Version
}

func getCmdVersion(cmdline string) string {
	cmd := exec.Command(cmdline, "-v")

	// 获取命令的输出，以及可能的错误信息
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "1.0.0"
	}

	return strings.TrimSuffix(string(output), "\n")
}

func GetCgiVersion() string {
	cmdline := SYSTEM_CONFIG_CGIPATH + CgiFile

	/*
		cmd := exec.Command(cmdline, "-v")

		// 获取命令的输出，以及可能的错误信息
		output, err := cmd.CombinedOutput()
		if err != nil {
			return "1.0.0"
		}

		return strings.TrimSuffix(string(output), "\n")
	*/
	return getCmdVersion(cmdline)

}

func GetRtuClientVersion() string {
	cmdline := SYSTEM_RUN_PATH + RtuFile
	return getCmdVersion(cmdline)
}
