package tools

import (
	"os/exec"
)

func RunNmap(target string) (string, error) {
	// 执行nmap扫描，输出XML格式
	cmd := exec.Command("nmap", "-Pn", "-sV", "-p 80", target)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
