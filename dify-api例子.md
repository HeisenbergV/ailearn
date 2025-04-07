# Dify 自定义工具
## 1. 实现 HTTP 服务

创建 Go 语言实现的 HTTP 服务：

```go
package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func sumHandler(c *gin.Context) {
	num1Str := c.Query("num1")
	num2Str := c.Query("num2")

	// 将字符串转换为数字
	num1, err1 := strconv.ParseFloat(num1Str, 64)
	num2, err2 := strconv.ParseFloat(num2Str, 64)

	// 检查转换是否成功
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请输入有效的数字",
		})
		return
	}

	// 计算和
	sum := num1 + num2

	c.JSON(http.StatusOK, sum)
}

func main() {
	r := gin.Default()

	// 设置 API 路由
	v1 := r.Group("/api/v1")
	{
		v1.GET("/sum", sumHandler)
	}

	// 启动服务器
	r.Run("0.0.0.0:8080")
}
```

## 2. 创建 OpenAPI 规范
OpenAPI 规范文件，定义 API 接口：
```json
{
  "openapi": "3.1.0",
  "info": {
    "title": "两数之和 API",
    "description": "计算两数之和",
    "version": "v1.0.0"
  },
  "servers": [
    {
      "url": "http://192.168.5.6:8080"
    }
  ],
  "paths": {
    "/api/v1/sum": {
      "get": {
        "description": "计算两数之和",
        "operationId": "sum",
        "parameters": [
          {
            "name": "num1",
            "in": "query",
            "description": "第一个数字",
            "required": true,
            "schema": {
              "type": "number"
            }
          },
          {
            "name": "num2",
            "in": "query",
            "description": "第二个数字",
            "required": true,
            "schema": {
              "type": "number"
            }
          }
        ],
        "deprecated": false
      }
    }
  },
  "components": {
    "schemas": {}
  }
}
```

## 3. Dify 集成步骤

1. 在 Dify 中创建 Agent 应用
2. 添加自定义工具（使用上述 OpenAPI 规范）
3. 启动 HTTP 服务
4. 在 Dify 中选择发布
5. 获取 API 地址（例如：http://192.168.5.7/chatbot/w0UEdifvJbehWOOy）

## 4. 使用示例

可以向 API 发送类似以下的问题：
- "1+2是多少"
- "22和1的和是多少"

系统会自动解析问题，调用自定义工具，并返回计算结果。
