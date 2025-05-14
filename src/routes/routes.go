package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "go-gin-test-job/docs"
	"go-gin-test-job/src/config"
	logger "go-gin-test-job/src/logger"
	middleware "go-gin-test-job/src/middlewares"
	accountModule "go-gin-test-job/src/modules/account"
	cronModule "go-gin-test-job/src/modules/cron"
	"strconv"
)

func New() (*gin.Engine, string) {
	initAppMode() // Set before gin.New()
	app := gin.Default()
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

	host := config.AppConfig.AppHost + ":" + strconv.Itoa(config.AppConfig.Port)
	return app, host
}

func initAppMode() {
	if config.AppConfig.IsDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}
