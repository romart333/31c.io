package testDatabase

import (
	"database/sql"
	"fmt"
	"go-gin-test-job/src/config"
	stringUtil "go-gin-test-job/src/utils/string"
	timeUtils "go-gin-test-job/src/utils/time"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"os"
	"path/filepath"
	"strings"
)

// Database instance
var DbConn *gorm.DB

func CreateDatabase(dbname string) error {
	var db *sql.DB
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?multiStatements=true", config.AppConfig.TestDatabase.Username, config.AppConfig.TestDatabase.Password, config.AppConfig.TestDatabase.Host, config.AppConfig.TestDatabase.Port)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s; CREATE DATABASE %s", dbname, dbname))
	if err != nil {
		return err
	}
	return nil
}

func DropDatabase(dbname string) {
	DbConn.Exec(fmt.Sprintf("DROP DATABASE %s", dbname))
}

func Connect() error {
	var err error
	err = CreateDatabase(config.AppConfig.TestDatabase.DbName)
	if err != nil {
		return err
	}
	// Use DSN string to open
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?multiStatements=true&parseTime=true", config.AppConfig.TestDatabase.Username, config.AppConfig.TestDatabase.Password, config.AppConfig.TestDatabase.Host, config.AppConfig.TestDatabase.Port, config.AppConfig.TestDatabase.DbName)
	DbConn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	err = DbConn.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{mysql.Open(dsn)},
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
	sqlDB.SetMaxIdleConns(config.AppConfig.TestDatabase.Connection.MaxNumber)
	sqlDB.SetMaxOpenConns(config.AppConfig.TestDatabase.Connection.OpenMaxNumber)
	sqlDB.SetConnMaxLifetime(timeUtils.DurationSeconds(config.AppConfig.TestDatabase.Connection.MaxLifetimeSec))
	return nil
}

func InitDatabase() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	sqlDir := filepath.Join(wd, "/scripts-mysql")
	files, err := filepath.Glob(filepath.Join(sqlDir, "*.sql"))
	if err != nil {
		return err
	}
	// Execute each SQL script
	for _, file := range files {
		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		query := string(sqlBytes)
		if stringUtil.CaseInsensitiveContains(query, "DELIMITER") {
			query = strings.ReplaceAll(query, "DELIMITER", "")
			query = strings.ReplaceAll(query, "$$", "")
		}
		DbConn.Exec(query)
		if err != nil {
			return fmt.Errorf("Execute sql script %s error. %s", file, err.Error())
		}
	}
	return nil
}
