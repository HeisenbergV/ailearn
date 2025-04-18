package svc

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"fogate/internal/config"
	"fogate/internal/types"
)

type ServiceContext struct {
	Config     config.Config
	ToolsCache *ToolsCache
	Terminal   *TerminalService
}

type ToolsCache struct {
	sync.RWMutex
	Total int
	Tools []types.SimpleToolInfo
}

func NewServiceContext(c config.Config) *ServiceContext {
	svcCtx := &ServiceContext{
		Config:     c,
		ToolsCache: &ToolsCache{},
		Terminal:   NewTerminalService(),
	}

	// 加载工具信息
	if err := svcCtx.loadTools(); err != nil {
		fmt.Printf("加载工具信息失败: %v\n", err)
	}

	return svcCtx
}

func (svc *ServiceContext) loadTools() error {
	toolsDir := "../tools"
	files, err := os.ReadDir(toolsDir)
	if err != nil {
		return fmt.Errorf("读取工具目录失败: %v", err)
	}

	var tools []types.SimpleToolInfo

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			filePath := filepath.Join(toolsDir, file.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("读取文件 %s 失败: %v\n", filePath, err)
				continue
			}

			var fullToolInfo types.ToolInfo
			if err := json.Unmarshal(content, &fullToolInfo); err != nil {
				fmt.Printf("解析文件 %s 失败: %v\n", filePath, err)
				continue
			}

			// 转换为简化版的工具信息
			simpleToolInfo := types.SimpleToolInfo{
				Name:        fullToolInfo.Tool.Name,
				Description: fullToolInfo.Tool.Description,
			}

			// 转换能力信息
			for _, cap := range fullToolInfo.Capabilities {
				simpleCap := types.SimpleCapability{
					Name:        cap.Name,
					Description: cap.Description,
				}
				simpleToolInfo.Capabilities = append(simpleToolInfo.Capabilities, simpleCap)
			}

			tools = append(tools, simpleToolInfo)
			fmt.Printf("成功加载工具: %s\n", simpleToolInfo.Name)
		}
	}

	// 更新缓存
	svc.ToolsCache.Lock()
	svc.ToolsCache.Tools = tools
	svc.ToolsCache.Total = len(tools)
	svc.ToolsCache.Unlock()

	fmt.Printf("总共加载了 %d 个工具\n", len(tools))
	return nil
}
