package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *handler) WelCome() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "This is SXX Micro Gateway!",
			"version": "Powered By Gin~",
		})
	}
}
