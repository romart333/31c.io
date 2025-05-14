package database

import (
	"database/sql"
	"go-gin-test-job/src/config"
	timeUtils "go-gin-test-job/src/utils/time"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

// Database instance
var DbConn *gorm.DB

var DefaultTxOptions = &sql.TxOptions{
	Isolation: sql.LevelReadCommitted,
	ReadOnly:  false,
}

func Connect() error {
	var err error
	DbConn, err = gorm.Open(mysql.Open(config.AppConfig.Database.Dsn), &gorm.Config{
		Logger: getDbLogger(),
	})
	if err != nil {
		return err
	}
	err = DbConn.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{mysql.Open(config.AppConfig.Database.Dsn)},
		Replicas: []gorm.Dialector{},
		// sources/replicas load balancing policy
		Policy: dbresolver.RandomPolicy{},
		// print sources/replicas mode in logger
		TraceResolverMode: true,
	}))
	if err != nil {
		return err
	}
	sqlDB, err := DbConn.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(config.AppConfig.Database.Connection.MaxNumber)
	sqlDB.SetMaxOpenConns(config.AppConfig.Database.Connection.OpenMaxNumber)
	sqlDB.SetConnMaxLifetime(timeUtils.DurationSeconds(config.AppConfig.Database.Connection.MaxLifetimeSec))
	return nil
}

func getDbLogger() logger.Interface {
	var dbLogger logger.Interface
	if config.AppConfig.Database.Logging {
		dbLogger = logger.Default.LogMode(logger.Info)
	} else {
		dbLogger = logger.Default.LogMode(logger.Silent)
	}
	return dbLogger
}
