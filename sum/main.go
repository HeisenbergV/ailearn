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
