package main

import (
	"net/http"
	"tool/tools"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Nikto API
	r.GET("/api/nikto", func(c *gin.Context) {
		target := c.Query("target")
		if target == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "target parameter is required"})
			return
		}
		data, _ := tools.RunNikto(target)
		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	})

	// Nmap API
	r.GET("/api/nmap", func(c *gin.Context) {
		target := c.Query("target")
		if target == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "target parameter is required"})
			return
		}
		data, _ := tools.RunNmap(target)
		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	})

	// SQLMap API
	r.GET("/api/sqlmap", func(c *gin.Context) {
		target := c.Query("target")
		if target == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "target parameter is required"})
			return
		}
		data, _ := tools.RunSqlmap(target)
		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	})

	r.Run(":8080")
}
