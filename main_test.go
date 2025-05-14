package main

import (
	"go-gin-test-job/src/config"
	"go-gin-test-job/test"
	testDatabase "go-gin-test-job/test/database"
	accountTests "go-gin-test-job/test/tests/account"
	cronTests "go-gin-test-job/test/tests/cron"
	"testing"
)

// TestMain runs before and after all test cases
func TestMain(m *testing.M) {
	test.InitApp()
	defer testDatabase.DropDatabase(config.AppConfig.TestDatabase.DbName)

	// Running integration
	m.Run()
}

func TestAllRoutes(t *testing.T) {
	t.Run("TestAccountRoute", accountTests.TestAccountRoute)
	t.Run("TestCronRoute", cronTests.TestCronRoute)
}
