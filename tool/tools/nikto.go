package tools

import (
	"os/exec"
)

func RunNikto(target string) (string, error) {
	// 执行nikto扫描，使用基本参数
	cmd := exec.Command("nikto", "-h", target)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
