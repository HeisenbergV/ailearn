package svc

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

type TerminalService struct {
	mu       sync.Mutex
	commands map[string]func(args []string) string
}

func NewTerminalService() *TerminalService {
	ts := &TerminalService{
		commands: make(map[string]func(args []string) string),
	}

	// 注册默认命令
	ts.RegisterCommand("help", ts.helpCommand)
	ts.RegisterCommand("tools", ts.toolsCommand)

	return ts
}

func (ts *TerminalService) RegisterCommand(name string, handler func(args []string) string) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.commands[name] = handler
}

func (ts *TerminalService) Start() {
	fmt.Println("欢迎使用 Fogate 终端")
	fmt.Println("输入 'help' 获取帮助信息")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("fogate> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		command := parts[0]
		args := parts[1:]

		ts.mu.Lock()
		handler, exists := ts.commands[command]
		ts.mu.Unlock()

		if !exists {
			fmt.Println("未知命令，输入 'help' 获取帮助信息")
			continue
		}

		result := handler(args)
		fmt.Println(result)
	}
}

// 默认命令处理函数
func (ts *TerminalService) helpCommand(args []string) string {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	var helpText strings.Builder
	helpText.WriteString("可用命令:\n")
	for cmd := range ts.commands {
		helpText.WriteString(fmt.Sprintf("  %s\n", cmd))
	}
	return helpText.String()
}

func (ts *TerminalService) toolsCommand(args []string) string {
	// 这里可以调用 ToolsCache 获取工具信息
	// 示例实现
	return "工具列表命令已执行"
}
