package main

import (
	"go-gin-test-job/src/config"
	"go-gin-test-job/src/database"
	"go-gin-test-job/src/logger"
	"go-gin-test-job/src/routes"
)

func init() {
	logger.InitializeLogger()
}

// @title Server API
// @version 1.0
// @description Server API

// @host localhost:3000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X_API_KEY
func main() {
	config.LoadConfig()
	if config.AppConfig.IsDebug {
		logger.SetDebugLevel()
	}
	if err := database.Connect(); err != nil {
		logger.Logger.Fatal().Msg("Connect to database error. Error - " + err.Error())
	}
	app, listenAddress := routes.New()
	if err := app.Run(listenAddress); err != nil {
		logger.Logger.Fatal().Msg("Startup error. Error - " + err.Error())
	}
}
