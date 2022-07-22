package services

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func RestartSystem() (err error) {
	cmd := exec.Command("sudo", "systemctl", "restart", "coco-captive-portal")
	err = cmd.Run()
	return
}

func GetNetInterfaceBytes(intf string) (rx, tx string) {
	cmd := exec.Command("cat", fmt.Sprintf("/sys/class/net/%s/statistics/rx_bytes", intf))
	bytes, err := cmd.Output()
	if err != nil {
		return
	}

	sBytes := string(bytes)
	sBytes = strings.Trim(sBytes, "\n")
	bInInt64, err := strconv.ParseInt(sBytes, 10, 64)
	if err != nil {
		return
	}

	rx = fmt.Sprintf("%.2f", float64(bInInt64)/1048576)

	cmd = exec.Command("cat", fmt.Sprintf("/sys/class/net/%s/statistics/tx_bytes", intf))
	bytes, err = cmd.Output()
	if err != nil {
		return
	}

	sBytes = string(bytes)
	sBytes = strings.Trim(sBytes, "\n")
	bInInt64, err = strconv.ParseInt(sBytes, 10, 64)
	if err != nil {
		return
	}

	tx = fmt.Sprintf("%.2f", float64(bInInt64)/1048576)

	return
}
