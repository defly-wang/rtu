package config

import (
	"fmt"
	"os/exec"

	"github.com/beevik/ntp"
)

const NtpServer = "pool.ntp.org"

func NtpDate() bool {
	rep, err := ntp.Query(NtpServer)
	//fmt.Println(rep.Time)
	if err == nil {
		now := rep.Time.Local()
		dateStr := fmt.Sprintf("%02d.%02d.%02d-%02d:%02d:%02d",
			now.Year(), now.Local().Month(), now.Day(),
			now.Hour(), now.Minute(), now.Second())
		//fmt.Println(dateStr)

		cmd := exec.Command("date", "-s", dateStr)
		cmd2 := exec.Command("hwclock", "-u", "-w")
		_, err := cmd.CombinedOutput()
		_, err2 := cmd2.CombinedOutput()

		return err == nil && err2 == nil

	}

	return false
}
