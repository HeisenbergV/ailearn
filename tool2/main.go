package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

type ToolInfo struct {
	Name             string          `json:"name"`
	ApplicableAssets []AssetInfo     `json:"applicable_assets"`
	Parameters       []ParameterInfo `json:"parameters"`
	Output           OutputInfo      `json:"output"`
}

type AssetInfo struct {
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

type ParameterInfo struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type OutputInfo struct {
	Format        string        `json:"format"`
	Normalization Normalization `json:"normalization"`
}

type Normalization struct {
	Type   string       `json:"type"`
	Schema []SchemaInfo `json:"schema"`
}

type SchemaInfo struct {
	Name        string       `json:"name"`
	Type        string       `json:"type"`
	Description string       `json:"description"`
	Items       []SchemaInfo `json:"items,omitempty"`
}

func main() {
	r := gin.Default()

	r.GET("/tools", func(c *gin.Context) {
		tools, err := readYAMLFiles()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, tools)
	})

	r.Run("0.0.0.0:8080")
}

func readYAMLFiles() ([]ToolInfo, error) {
	var tools []ToolInfo

	files, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".yml") {
			continue
		}

		data, err := os.ReadFile(file.Name())
		if err != nil {
			log.Printf("Error reading file %s: %v", file.Name(), err)
			continue
		}

		var tool ToolInfo
		if err := yaml.Unmarshal(data, &tool); err != nil {
			log.Printf("Error parsing YAML in %s: %v", file.Name(), err)
			continue
		}

		// 填充null字段
		if tool.ApplicableAssets == nil {
			tool.ApplicableAssets = []AssetInfo{
				{
					Type:     "ip",
					Required: true,
				},
			}
		}

		tools = append(tools, tool)
	}

	return tools, nil
}
