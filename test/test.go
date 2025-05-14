package test

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go-gin-test-job/src/config"
	appDatabase "go-gin-test-job/src/database"
	"go-gin-test-job/src/database/entities"
	"go-gin-test-job/src/logger"
	accountModuleDto "go-gin-test-job/src/modules/account/dto"
	testDatabase "go-gin-test-job/test/database"
	testRoutes "go-gin-test-job/test/routes"
	"go-gin-test-job/test/seeds"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

var TestApp *gin.Engine
var TestAppConfig *TestServerConfig

type TestServerConfig struct {
	Port int
	Host string
}

type ErrorResponseDto struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (config *TestServerConfig) GetUrl() string {
	return TestAppConfig.Host + ":" + strconv.Itoa(TestAppConfig.Port)
}

func init() {
	logger.InitializeLogger()
}

func InitApp() *gin.Engine {
	config.LoadConfig()
	// Connect to databases
	if err := testDatabase.Connect(); err != nil {
		logger.Logger.Fatal().Msg("Connect to database error. Error - " + err.Error())
	}
	// Init database
	if err := testDatabase.InitDatabase(); err != nil {
		logger.Logger.Fatal().Msg("Init database error. Error - " + err.Error())
	}
	// Insert data to database
	createTestData()
	// Set DbConn for app
	appDatabase.DbConn = testDatabase.DbConn
	TestAppConfig = &TestServerConfig{
		Host: "localhost",
		Port: 8080,
	}
	TestApp = testRoutes.New()
	// Start test server
	go TestApp.Run(TestAppConfig.GetUrl())
	// Give server some time to start
	time.Sleep(100 * time.Millisecond)
	return TestApp
}

func createTestData() {
	// Add accounts
	for _, account := range seeds.FillAccountList() {
		testDatabase.DbConn.Create(&account)
	}
}

func TestListSort[T any](list []T, orderBy string) bool {
	if len(list) == 0 {
		return true
	}
	sortParams := strings.Split(orderBy, " ")
	_, fieldType, exists := getFieldInfo(list[0], sortParams[0])
	if !exists {
		return false
	}
	isSorted := true
	for i := 1; i < len(list); i++ {
		before, exists1 := getFieldValue(list[i-1], sortParams[0])
		current, exists2 := getFieldValue(list[i], sortParams[0])
		if !exists1 || !exists2 {
			return false
		}
		// Type assertion
		if fieldType.Kind() == reflect.Int {
			beforeValue, ok1 := before.(int)
			currentValue, ok2 := current.(int)
			// Ensure both values are valid before comparison
			if !ok1 || !ok2 {
				return false
			}
			if sortParams[1] == "ASC" && beforeValue > currentValue || sortParams[1] == "DESC" && beforeValue < currentValue {
				isSorted = false
				break
			}
		} else if fieldType.Kind() == reflect.Int64 {
			beforeValue, ok1 := before.(int64)
			currentValue, ok2 := current.(int64)
			// Ensure both values are valid before comparison
			if !ok1 || !ok2 {
				return false
			}
			if sortParams[1] == "ASC" && beforeValue > currentValue || sortParams[1] == "DESC" && beforeValue < currentValue {
				isSorted = false
				break
			}
		} else if fieldType.Kind() == reflect.String {
			beforeValue, ok1 := before.(string)
			currentValue, ok2 := current.(string)
			// Ensure both values are valid before comparison
			if !ok1 || !ok2 {
				return false
			}
			if sortParams[1] == "ASC" && beforeValue > currentValue || sortParams[1] == "DESC" && beforeValue < currentValue {
				isSorted = false
				break
			}
		} else {
			logger.Logger.Fatal().Msg("Unsupported type")
		}

	}
	return isSorted
}

func getFieldInfo(obj interface{}, jsonFieldName string) (interface{}, reflect.Type, bool) {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	// Loop through struct fields
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")

		// Match JSON field name
		if jsonTag == jsonFieldName {
			return v.Field(i).Interface(), field.Type, true
		}
	}
	return nil, nil, false
}

func getFieldValue(obj interface{}, fieldName string) (interface{}, bool) {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	// Loop through struct fields to find the JSON tag
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == fieldName {
			return v.Field(i).Interface(), true
		}
	}
	return nil, false
}

func CompareAccount(t *testing.T, account *entities.Account, accountDto accountModuleDto.AccountDto) {
	assert.Equal(t, account.Id, accountDto.Id)
	assert.Equal(t, account.Address, accountDto.Address)
	assert.Equal(t, string(account.Status), accountDto.Status)
	assert.Equal(t, account.CreatedAt, accountDto.CreatedAt)
	assert.Equal(t, account.UpdatedAt, accountDto.UpdatedAt)
}
