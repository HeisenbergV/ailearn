package tools

import (
	"os/exec"
)

func RunSqlmap(target string) (string, error) {
	// 执行sqlmap扫描，使用POST请求和表单数据
	cmd := exec.Command("sqlmap",
		"-u", target,
		"--data", "username=admin&password=password&Login=Login",
		"--batch",
		"--random-agent")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
