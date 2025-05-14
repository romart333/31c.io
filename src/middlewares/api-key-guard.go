package middleware

import (
	"github.com/gin-gonic/gin"
	errorHelper "go-gin-test-job/src/common/error-helpers"
	"go-gin-test-job/src/config"
)

func AdminApiKeyGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" || config.AppConfig.AdminXApiKey == "" || apiKey != config.AppConfig.AdminXApiKey {
			_ = errorHelper.RespondUnauthorizedError(c)
			c.Abort()
			return
		}
		c.Next()
	}
}

func CronApiKeyGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" || config.AppConfig.CronXApiKey == "" || apiKey != config.AppConfig.CronXApiKey {
			_ = errorHelper.RespondUnauthorizedError(c)
			c.Abort()
			return
		}
		c.Next()
	}
}
