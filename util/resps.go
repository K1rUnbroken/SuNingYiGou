package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// RespSuccess 正确返回
func RespSuccess(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": msg,
		"data":    data,
	})
}

// RespClientError 客户端出错
func RespClientError(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": msg,
	})
}
