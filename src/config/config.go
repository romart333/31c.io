package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"go-gin-test-job/src/logger"
	typeUtil "go-gin-test-job/src/utils/type"
	"os"
	"strconv"
)

type DbConnectionConfig struct {
	MaxNumber      int
	OpenMaxNumber  int
	MaxLifetimeSec int
}

type DbConfig struct {
	Dsn        string
	Connection DbConnectionConfig
	Logging    bool
}

type TestDbConfig struct {
	Host       string
	Port       int
	Username   string
	Password   string
	DbName     string
	Connection DbConnectionConfig
	Logging    bool
}

type Config struct {
	AppName           string
	AppHost           string
	Port              int
	IsDebug           bool
	AdminXApiKey      string
	CronXApiKey       string
	RequestTimeoutSec int
	CronBatchCount    int
	Database          DbConfig
	TestDatabase      TestDbConfig
}

var AppConfig *Config

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		logger.Logger.Fatal().Msg("Loading .env file error. Error - " + err.Error())
	}
	appName := getEnvAsString("APP_NAME", typeUtil.String("TestApp"))
	appHost := getEnvAsString("APP_HOST", typeUtil.String("localhost"))
	port := getEnvAsInt("PORT", typeUtil.Int(3000))
	isDebug := getEnvAsBool("IS_DEBUG", typeUtil.Bool(true))
	adminXApiKey := getEnvAsString("ADMIN_X_API_KEY", nil)
	cronXApiKey := getEnvAsString("CRON_X_API_KEY", nil)
	requestTimeoutSec := getEnvAsInt("REQUEST_TIMEOUT_SEC", typeUtil.Int(20))
	cronBatchCount := getEnvAsInt("CRON_BATCH_COUNT", typeUtil.Int(5))

	dbHost := getEnvAsString("DB_HOST", typeUtil.String("localhost"))
	dbPort := getEnvAsInt("DB_PORT", typeUtil.Int(3306))
	dbUsername := getEnvAsString("DB_USERNAME", typeUtil.String("username"))
	dbPassword := getEnvAsString("DB_PASSWORD", typeUtil.String("password"))
	dbSchema := getEnvAsString("DB_SCHEMA", typeUtil.String("database"))
	dbDns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbUsername, dbPassword, dbHost, dbPort, dbSchema)

	testDbHost := getEnvAsString("TEST_DB_HOST", typeUtil.String("localhost"))
	testDbPort := getEnvAsInt("TEST_DB_PORT", typeUtil.Int(3406))
	testDbUsername := getEnvAsString("TEST_DB_USERNAME", typeUtil.String("root"))
	testDbPassword := getEnvAsString("TEST_DB_PASSWORD", typeUtil.String("root_password"))
	testDbSchema := getEnvAsString("TEST_DB_SCHEMA", typeUtil.String("server"))

	defaultDbConnection := DbConnectionConfig{
		MaxNumber:      10,
		OpenMaxNumber:  100,
		MaxLifetimeSec: 3600,
	}

	AppConfig = &Config{
		AppName:           appName,
		AppHost:           appHost,
		Port:              port,
		IsDebug:           isDebug,
		AdminXApiKey:      adminXApiKey,
		CronXApiKey:       cronXApiKey,
		RequestTimeoutSec: requestTimeoutSec,
		CronBatchCount:    cronBatchCount,
		Database: DbConfig{
			Dsn:        dbDns,
			Connection: defaultDbConnection,
			Logging:    false,
		},
		TestDatabase: TestDbConfig{
			Host:       testDbHost,
			Port:       testDbPort,
			Username:   testDbUsername,
			Password:   testDbPassword,
			DbName:     testDbSchema,
			Connection: defaultDbConnection,
			Logging:    false,
		},
	}
}

func getEnvAsString(key string, defaultValue *string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		if defaultValue == nil {
			logger.Logger.Fatal().Msg(fmt.Sprintf("Required environment variable %s is not set", key))
		}
		return *defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue *int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		if defaultValue == nil {
			logger.Logger.Fatal().Msg(fmt.Sprintf("Required environment variable %s is not set", key))
		}
		return *defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		logger.Logger.Fatal().Msg(fmt.Sprintf("Environment variable %s must be an integer, got %s", key, value))
	}
	return intValue
}

func getEnvAsBool(key string, defaultValue *bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		if defaultValue == nil {
			logger.Logger.Fatal().Msg(fmt.Sprintf("Required environment variable %s is not set", key))
		}
		return *defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		logger.Logger.Fatal().Msg(fmt.Sprintf("Environment variable %s must be a valid boolean, got %s", key, value))
	}
	return boolValue
}
