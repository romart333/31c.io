package testRoutes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	logger "go-gin-test-job/src/logger"
	middleware "go-gin-test-job/src/middlewares"
	accountModule "go-gin-test-job/src/modules/account"
	cronModule "go-gin-test-job/src/modules/cron"
)

func New() *gin.Engine {
	initAppMode()
	app := gin.New()
	_ = app.SetTrustedProxies(nil)

	// Set up middleware
	app.Use(gin.Recovery())
	app.Use(cors.Default())
	app.Use(logger.LogMiddleware())
	app.Use(middleware.RequestIDMiddleware())
	app.Use(middleware.ErrorHandler())

	// Swagger handler
	app.GET("/api/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Account routes
	accountMethods := app.Group("/account")
	accountMethods.GET("", middleware.AdminApiKeyGuard(), accountModule.GetAccounts)
	accountMethods.POST("", middleware.AdminApiKeyGuard(), accountModule.CreateAccount)

	// Cron routes
	cronMethods := app.Group("/cron")
	cronMethods.POST("/account-balance", middleware.CronApiKeyGuard(), cronModule.UpdateAccountsBalances)

	return app
}

func initAppMode() {
	gin.SetMode(gin.TestMode)
}
