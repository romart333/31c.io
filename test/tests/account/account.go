package accountTests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	errorHelpers "go-gin-test-job/src/common/error-helpers"
	"go-gin-test-job/src/config"
	"go-gin-test-job/src/database"
	"go-gin-test-job/src/database/entities"
	accountModuleDto "go-gin-test-job/src/modules/account/dto"
	arrayUtil "go-gin-test-job/src/utils/array"
	numberUtil "go-gin-test-job/src/utils/number"
	orderUtil "go-gin-test-job/src/utils/order"
	timeUtil "go-gin-test-job/src/utils/time"
	"go-gin-test-job/test"
	"go-gin-test-job/test/seeds"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestAccountRoute(t *testing.T) {
	// GetAccounts
	validationGetAccountsTests(t)
	t.Run("TestGetAccountsRoute_SuccessNoParams", TestGetAccountsRoute_SuccessNoParams)
	t.Run("TestGetAccountsRoute_SuccessParamsOffsetAndCount", TestGetAccountsRoute_SuccessParamsOffsetAndCount)
	t.Run("TestGetAccountsRoute_SuccessParamsStatus", TestGetAccountsRoute_SuccessParamsStatus)
	t.Run("TestGetAccountsRoute_SuccessParamsOrderBy", TestGetAccountsRoute_SuccessParamsOrderBy)
	t.Run("TestGetAccountsRoute_SuccessParamsStatusAndOrderBy", TestGetAccountsRoute_SuccessParamsStatusAndOrderBy)
	t.Run("TestGetAccountsRoute_SuccessParamsOffsetAndCountAndStatusAndOrderBy", TestGetAccountsRoute_SuccessParamsOffsetAndCountAndStatusAndOrderBy)
	// CreateAccount
	validationCreateAccountTests(t)
	t.Run("TestCreateAccountRoute_FailAddressAlreadyExists", TestCreateAccountRoute_FailAddressAlreadyExists)
	t.Run("TestCreateAccountRoute_Success", TestCreateAccountRoute_Success)
}

func validationGetAccountsTests(t *testing.T) {
	validationTests := []struct {
		name         string
		params       accountModuleDto.GetAccountRequestDto
		expectedCode int
		expectedBody errorHelpers.ResponseBadRequestErrorHTTP
	}{
		{
			"FailInvalidOffsetMinValue",
			accountModuleDto.GetAccountRequestDto{Offset: -5},
			http.StatusBadRequest,
			errorHelpers.ResponseBadRequestErrorHTTP{Success: false, Message: "Offset must be greater than or equal 0"},
		},
		{
			"FailInvalidCountMinValue",
			accountModuleDto.GetAccountRequestDto{Count: -1},
			http.StatusBadRequest,
			errorHelpers.ResponseBadRequestErrorHTTP{Success: false, Message: "Count must be greater than or equal 1"},
		},
		{
			"FailInvalidCountMaxValue",
			accountModuleDto.GetAccountRequestDto{Count: 101},
			http.StatusBadRequest,
			errorHelpers.ResponseBadRequestErrorHTTP{Success: false, Message: "Count must be less than or equal 100"},
		},
		{
			"FailInvalidStatus",
			accountModuleDto.GetAccountRequestDto{Status: "invalid status"},
			http.StatusBadRequest,
			errorHelpers.ResponseBadRequestErrorHTTP{Success: false, Message: fmt.Sprintf("%s must be one of the next values: %s", "Status", strings.Join(entities.AccountStatusList, ","))},
		},
		{
			"FailInvalidOrderBy",
			accountModuleDto.GetAccountRequestDto{OrderBy: "invalid order by"},
			http.StatusBadRequest,
			errorHelpers.ResponseBadRequestErrorHTTP{Success: false, Message: "invalid order by parameter: invalid order by"},
		},
		{
			"FailInvalidOrderByMaxLength",
			accountModuleDto.GetAccountRequestDto{OrderBy: strings.Repeat("OrderBy", 255)},
			http.StatusBadRequest,
			errorHelpers.ResponseBadRequestErrorHTTP{Success: false, Message: "OrderBy must be shorter than or equal to 255 characters"},
		},
	}
	for _, validationTest := range validationTests {
		t.Run("TestGetAccountsRoute"+validationTest.name, func(t *testing.T) {
			type Params struct {
				Count   int                    `json:"count"`
				Offset  int                    `json:"offset"`
				Status  entities.AccountStatus `json:"status"`
				OrderBy string                 `json:"orderBy"`
			}
			params := &Params{
				Count:   validationTest.params.Count,
				Offset:  validationTest.params.Offset,
				Status:  validationTest.params.Status,
				OrderBy: validationTest.params.OrderBy,
			}

			query := url.Values{}
			query.Add("count", numberUtil.IntToString(params.Count))
			query.Add("offset", numberUtil.IntToString(params.Offset))
			query.Add("status", string(params.Status))
			query.Add("orderBy", params.OrderBy)

			u := &url.URL{
				Path:     fmt.Sprintf("/account"),
				RawQuery: query.Encode(),
			}

			response := httptest.NewRecorder()
			request := httptest.NewRequest("GET", u.String(), nil)
			request.Header.Set("X-API-Key", config.AppConfig.AdminXApiKey)
			test.TestApp.ServeHTTP(response, request)
			assert.Equal(t, validationTest.expectedCode, response.Code)

			// Read the response body and parse JSON
			var responseDto errorHelpers.ResponseBadRequestErrorHTTP
			err := json.NewDecoder(response.Body).Decode(&responseDto)
			assert.Nil(t, err)

			assert.NotNil(t, responseDto.Success, "Success parameter should exist")
			assert.NotNil(t, responseDto.Message, "Message parameter should exist")

			assert.Equal(t, validationTest.expectedBody.Success, responseDto.Success)
			assert.Equal(t, validationTest.expectedBody.Message, responseDto.Message)
		})
	}
}

func TestGetAccountsRoute_SuccessNoParams(t *testing.T) {
	u := &url.URL{
		Path: fmt.Sprintf("/account"),
	}

	accounts, total := database.GetAccountsAndTotal("", make(map[string]string), accountModuleDto.DEFAULT_ACCOUNT_OFFSET, accountModuleDto.DEFAULT_ACCOUNT_COUNT)

	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", u.String(), nil)
	request.Header.Set("X-API-Key", config.AppConfig.AdminXApiKey)
	test.TestApp.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code)

	// Read the response body and parse JSON
	var responseDto accountModuleDto.GetAccountResponseDto
	err := json.NewDecoder(response.Body).Decode(&responseDto)
	assert.Nil(t, err)

	assert.NotNil(t, responseDto.Offset, "Offset parameter should exist")
	assert.NotNil(t, responseDto.Count, "Count parameter should exist")
	assert.NotNil(t, responseDto.Total, "Total parameter should exist")
	assert.NotNil(t, responseDto.List, "List parameter should exist")

	assert.Equal(t, 0, responseDto.Offset)
	assert.Equal(t, accountModuleDto.DEFAULT_ACCOUNT_COUNT, responseDto.Count)
	assert.Equal(t, total, responseDto.Total)
	assert.Equal(t, len(accounts), len(responseDto.List))

	for _, accountDto := range responseDto.List {
		conditions := []func(account *entities.Account) bool{
			func(a *entities.Account) bool {
				return a.Id == accountDto.Id
			},
		}
		account := arrayUtil.FindItem(accounts, conditions)
		test.CompareAccount(t, *account, accountDto)
	}
}

func TestGetAccountsRoute_SuccessParamsOffsetAndCount(t *testing.T) {
	type Params struct {
		Count  int `json:"count"`
		Offset int `json:"offset"`
	}
	params := &Params{
		Count:  2,
		Offset: 1,
	}

	query := url.Values{}
	query.Add("count", numberUtil.IntToString(params.Count))
	query.Add("offset", numberUtil.IntToString(params.Offset))

	u := &url.URL{
		Path:     fmt.Sprintf("/account"),
		RawQuery: query.Encode(),
	}

	accounts, total := database.GetAccountsAndTotal("", make(map[string]string), params.Offset, params.Count)

	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", u.String(), nil)
	request.Header.Set("X-API-Key", config.AppConfig.AdminXApiKey)
	test.TestApp.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code)

	// Read the response body and parse JSON
	var responseDto accountModuleDto.GetAccountResponseDto
	err := json.NewDecoder(response.Body).Decode(&responseDto)
	assert.Nil(t, err)

	assert.NotNil(t, responseDto.Offset, "Offset parameter should exist")
	assert.NotNil(t, responseDto.Count, "Count parameter should exist")
	assert.NotNil(t, responseDto.Total, "Total parameter should exist")
	assert.NotNil(t, responseDto.List, "List parameter should exist")

	assert.Equal(t, params.Offset, responseDto.Offset)
	assert.Equal(t, params.Count, responseDto.Count)
	assert.Equal(t, total, responseDto.Total)
	assert.Equal(t, len(accounts), len(responseDto.List))

	for _, accountDto := range responseDto.List {
		conditions := []func(account *entities.Account) bool{
			func(a *entities.Account) bool {
				return a.Id == accountDto.Id
			},
		}
		account := arrayUtil.FindItem(accounts, conditions)
		test.CompareAccount(t, *account, accountDto)
	}
}

func TestGetAccountsRoute_SuccessParamsStatus(t *testing.T) {
	type Params struct {
		Status entities.AccountStatus `json:"status"`
	}
	params := &Params{
		Status: entities.AccountStatusOn,
	}

	query := url.Values{}
	query.Add("status", string(params.Status))

	u := &url.URL{
		Path:     fmt.Sprintf("/account"),
		RawQuery: query.Encode(),
	}

	accounts, total := database.GetAccountsAndTotal(params.Status, make(map[string]string), accountModuleDto.DEFAULT_ACCOUNT_OFFSET, accountModuleDto.DEFAULT_ACCOUNT_COUNT)

	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", u.String(), nil)
	request.Header.Set("X-API-Key", config.AppConfig.AdminXApiKey)
	test.TestApp.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code)

	// Read the response body and parse JSON
	var responseDto accountModuleDto.GetAccountResponseDto
	err := json.NewDecoder(response.Body).Decode(&responseDto)
	assert.Nil(t, err)

	assert.NotNil(t, responseDto.Offset, "Offset parameter should exist")
	assert.NotNil(t, responseDto.Count, "Count parameter should exist")
	assert.NotNil(t, responseDto.Total, "Total parameter should exist")
	assert.NotNil(t, responseDto.List, "List parameter should exist")

	assert.Equal(t, accountModuleDto.DEFAULT_ACCOUNT_OFFSET, responseDto.Offset)
	assert.Equal(t, accountModuleDto.DEFAULT_ACCOUNT_COUNT, responseDto.Count)
	assert.Equal(t, total, responseDto.Total)
	assert.Equal(t, len(accounts), len(responseDto.List))

	for _, accountDto := range responseDto.List {
		conditions := []func(account *entities.Account) bool{
			func(a *entities.Account) bool {
				return a.Id == accountDto.Id
			},
		}
		account := arrayUtil.FindItem(accounts, conditions)
		test.CompareAccount(t, *account, accountDto)
	}
}

func TestGetAccountsRoute_SuccessParamsOrderBy(t *testing.T) {
	type Params struct {
		OrderBy string `json:"orderBy"`
	}
	params := &Params{
		OrderBy: "id DESC",
	}

	query := url.Values{}
	query.Add("orderBy", params.OrderBy)

	u := &url.URL{
		Path:     fmt.Sprintf("/account"),
		RawQuery: query.Encode(),
	}

	orderParams, err := orderUtil.GetOrderByParamsSecure(nil, params.OrderBy, ",", accountModuleDto.GetAvailableAccountSortFieldList)
	accounts, total := database.GetAccountsAndTotal("", orderParams, accountModuleDto.DEFAULT_ACCOUNT_OFFSET, accountModuleDto.DEFAULT_ACCOUNT_COUNT)

	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", u.String(), nil)
	request.Header.Set("X-API-Key", config.AppConfig.AdminXApiKey)
	test.TestApp.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code)

	// Read the response body and parse JSON
	var responseDto accountModuleDto.GetAccountResponseDto
	err = json.NewDecoder(response.Body).Decode(&responseDto)
	assert.Nil(t, err)

	assert.NotNil(t, responseDto.Offset, "Offset parameter should exist")
	assert.NotNil(t, responseDto.Count, "Count parameter should exist")
	assert.NotNil(t, responseDto.Total, "Total parameter should exist")
	assert.NotNil(t, responseDto.List, "List parameter should exist")

	assert.Equal(t, accountModuleDto.DEFAULT_ACCOUNT_OFFSET, responseDto.Offset)
	assert.Equal(t, accountModuleDto.DEFAULT_ACCOUNT_COUNT, responseDto.Count)
	assert.Equal(t, total, responseDto.Total)
	assert.Equal(t, len(accounts), len(responseDto.List))

	for _, accountDto := range responseDto.List {
		conditions := []func(account *entities.Account) bool{
			func(a *entities.Account) bool {
				return a.Id == accountDto.Id
			},
		}
		account := arrayUtil.FindItem(accounts, conditions)
		test.CompareAccount(t, *account, accountDto)
	}

	assert.Equal(t, true, test.TestListSort(responseDto.List, params.OrderBy), "List is not sorted")
}

func TestGetAccountsRoute_SuccessParamsStatusAndOrderBy(t *testing.T) {
	type Params struct {
		Status  entities.AccountStatus `json:"status"`
		OrderBy string                 `json:"orderBy"`
	}
	params := &Params{
		Status:  entities.AccountStatusOff,
		OrderBy: "updated_at DESC",
	}

	query := url.Values{}
	query.Add("status", string(params.Status))
	query.Add("orderBy", params.OrderBy)

	u := &url.URL{
		Path:     fmt.Sprintf("/account"),
		RawQuery: query.Encode(),
	}

	orderParams, err := orderUtil.GetOrderByParamsSecure(nil, params.OrderBy, ",", accountModuleDto.GetAvailableAccountSortFieldList)
	accounts, total := database.GetAccountsAndTotal(params.Status, orderParams, accountModuleDto.DEFAULT_ACCOUNT_OFFSET, accountModuleDto.DEFAULT_ACCOUNT_COUNT)

	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", u.String(), nil)
	request.Header.Set("X-API-Key", config.AppConfig.AdminXApiKey)
	test.TestApp.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code)

	// Read the response body and parse JSON
	var responseDto accountModuleDto.GetAccountResponseDto
	err = json.NewDecoder(response.Body).Decode(&responseDto)
	assert.Nil(t, err)

	assert.NotNil(t, responseDto.Offset, "Offset parameter should exist")
	assert.NotNil(t, responseDto.Count, "Count parameter should exist")
	assert.NotNil(t, responseDto.Total, "Total parameter should exist")
	assert.NotNil(t, responseDto.List, "List parameter should exist")

	assert.Equal(t, accountModuleDto.DEFAULT_ACCOUNT_OFFSET, responseDto.Offset)
	assert.Equal(t, accountModuleDto.DEFAULT_ACCOUNT_COUNT, responseDto.Count)
	assert.Equal(t, total, responseDto.Total)
	assert.Equal(t, len(accounts), len(responseDto.List))

	for _, accountDto := range responseDto.List {
		conditions := []func(account *entities.Account) bool{
			func(a *entities.Account) bool {
				return a.Id == accountDto.Id
			},
		}
		account := arrayUtil.FindItem(accounts, conditions)
		test.CompareAccount(t, *account, accountDto)
	}

	assert.Equal(t, true, test.TestListSort(responseDto.List, params.OrderBy), "List is not sorted")
}

func TestGetAccountsRoute_SuccessParamsOffsetAndCountAndStatusAndOrderBy(t *testing.T) {
	type Params struct {
		Count   int                    `json:"count"`
		Offset  int                    `json:"offset"`
		Status  entities.AccountStatus `json:"status"`
		OrderBy string                 `json:"orderBy"`
	}
	params := &Params{
		Count:   2,
		Offset:  0,
		Status:  entities.AccountStatusOff,
		OrderBy: "updated_at ASC",
	}

	query := url.Values{}
	query.Add("count", numberUtil.IntToString(params.Count))
	query.Add("offset", numberUtil.IntToString(params.Offset))
	query.Add("status", string(params.Status))
	query.Add("orderBy", params.OrderBy)

	u := &url.URL{
		Path:     fmt.Sprintf("/account"),
		RawQuery: query.Encode(),
	}

	orderParams, err := orderUtil.GetOrderByParamsSecure(nil, params.OrderBy, ",", accountModuleDto.GetAvailableAccountSortFieldList)
	accounts, total := database.GetAccountsAndTotal(params.Status, orderParams, params.Offset, params.Count)

	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", u.String(), nil)
	request.Header.Set("X-API-Key", config.AppConfig.AdminXApiKey)
	test.TestApp.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code)

	// Read the response body and parse JSON
	var responseDto accountModuleDto.GetAccountResponseDto
	err = json.NewDecoder(response.Body).Decode(&responseDto)
	assert.Nil(t, err)

	assert.NotNil(t, responseDto.Offset, "Offset parameter should exist")
	assert.NotNil(t, responseDto.Count, "Count parameter should exist")
	assert.NotNil(t, responseDto.Total, "Total parameter should exist")
	assert.NotNil(t, responseDto.List, "List parameter should exist")

	assert.Equal(t, params.Offset, responseDto.Offset)
	assert.Equal(t, params.Count, responseDto.Count)
	assert.Equal(t, total, responseDto.Total)
	assert.Equal(t, len(accounts), len(responseDto.List))

	for _, accountDto := range responseDto.List {
		conditions := []func(account *entities.Account) bool{
			func(a *entities.Account) bool {
				return a.Id == accountDto.Id
			},
		}
		account := arrayUtil.FindItem(accounts, conditions)
		test.CompareAccount(t, *account, accountDto)
	}

	assert.Equal(t, true, test.TestListSort(responseDto.List, params.OrderBy), "List is not sorted")
}

func validationCreateAccountTests(t *testing.T) {
	validationTests := []struct {
		name         string
		params       accountModuleDto.PostCreateAccountRequestDto
		expectedCode int
		expectedBody errorHelpers.ResponseBadRequestErrorHTTP
	}{
		{
			"FailNoBody",
			accountModuleDto.PostCreateAccountRequestDto{},
			http.StatusBadRequest,
			errorHelpers.ResponseBadRequestErrorHTTP{Success: false, Message: "Address format is wrong"},
		},
		{
			"FailInvalidAddress",
			accountModuleDto.PostCreateAccountRequestDto{Address: "invalid address"},
			http.StatusBadRequest,
			errorHelpers.ResponseBadRequestErrorHTTP{Success: false, Message: "Address format is wrong"},
		},
		{
			"FailInvalidStatus",
			accountModuleDto.PostCreateAccountRequestDto{Address: "14yqg2y3a6HMgW9MiF5tVPAH4Dr1uxGKFJ", Status: "invalid status"},
			http.StatusBadRequest,
			errorHelpers.ResponseBadRequestErrorHTTP{Success: false, Message: fmt.Sprintf("%s must be one of the next values: %s", "Status", strings.Join(entities.AccountStatusList, ","))},
		},
	}
	for _, validationTest := range validationTests {
		t.Run("TestCreateAccountRoute"+validationTest.name, func(t *testing.T) {
			type Params struct {
				Address string                 `json:"address"`
				Status  entities.AccountStatus `json:"status"`
			}
			params := &Params{
				Address: validationTest.params.Address,
				Status:  validationTest.params.Status,
			}
			body, _ := json.Marshal(params)

			u := &url.URL{
				Path: fmt.Sprintf("/account"),
			}

			response := httptest.NewRecorder()
			request := httptest.NewRequest("POST", u.String(), bytes.NewBuffer(body))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("X-API-Key", config.AppConfig.AdminXApiKey)
			test.TestApp.ServeHTTP(response, request)
			assert.Equal(t, validationTest.expectedCode, response.Code)

			// Read the response body and parse JSON
			var responseDto errorHelpers.ResponseBadRequestErrorHTTP
			err := json.NewDecoder(response.Body).Decode(&responseDto)
			assert.Nil(t, err)

			assert.NotNil(t, responseDto.Success, "Success parameter should exist")
			assert.NotNil(t, responseDto.Message, "Message parameter should exist")

			assert.Equal(t, validationTest.expectedBody.Success, responseDto.Success)
			assert.Equal(t, validationTest.expectedBody.Message, responseDto.Message)
		})
	}
}

func TestCreateAccountRoute_FailAddressAlreadyExists(t *testing.T) {
	accountInfo := seeds.ACCOUNTS.ACCOUNT_1
	type Params struct {
		Address string                 `json:"address"`
		Status  entities.AccountStatus `json:"status"`
	}
	params := &Params{
		Address: accountInfo.Address,
		Status:  entities.AccountStatusOn,
	}
	body, _ := json.Marshal(params)

	u := &url.URL{
		Path: fmt.Sprintf("/account"),
	}

	assert.Equal(t, true, database.IsAddressExists(nil, params.Address), "Address must exists")

	response := httptest.NewRecorder()
	request := httptest.NewRequest("POST", u.String(), bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-Key", config.AppConfig.AdminXApiKey)
	test.TestApp.ServeHTTP(response, request)
	assert.Equal(t, http.StatusConflict, response.Code)

	// Read the response body and parse JSON
	var responseDto errorHelpers.ResponseBadRequestErrorHTTP
	err := json.NewDecoder(response.Body).Decode(&responseDto)
	assert.Nil(t, err)

	assert.NotNil(t, responseDto.Success, "Success parameter should exist")
	assert.NotNil(t, responseDto.Message, "Message parameter should exist")

	assert.Equal(t, false, responseDto.Success)
	assert.Equal(t, "Address already exists", responseDto.Message)
}

func TestCreateAccountRoute_Success(t *testing.T) {
	start := timeUtil.GetUnixTime()
	type Params struct {
		Address string                 `json:"address"`
		Status  entities.AccountStatus `json:"status"`
	}
	params := &Params{
		Address: "32AaKxGbdhGMSGutcZjspFq9U89jJHW1um",
		Status:  entities.AccountStatusOn,
	}
	body, _ := json.Marshal(params)

	u := &url.URL{
		Path: fmt.Sprintf("/account"),
	}

	assert.Equal(t, false, database.IsAddressExists(nil, params.Address), "Address must not exists")

	response := httptest.NewRecorder()
	request := httptest.NewRequest("POST", u.String(), bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-Key", config.AppConfig.AdminXApiKey)
	test.TestApp.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code)

	// Read the response body and parse JSON
	var responseDto accountModuleDto.AccountDto
	err := json.NewDecoder(response.Body).Decode(&responseDto)
	assert.Nil(t, err)

	assert.NotNil(t, responseDto.Id, "Id parameter should exist")
	assert.NotNil(t, responseDto.Address, "Address parameter should exist")
	assert.NotNil(t, responseDto.Balance, "Balance parameter should exist")
	assert.NotNil(t, responseDto.Status, "Status parameter should exist")
	assert.NotNil(t, responseDto.CreatedAt, "CreatedAt parameter should exist")
	assert.NotNil(t, responseDto.UpdatedAt, "UpdatedAt parameter should exist")

	accountAfter := database.GetAccountByAddress(params.Address)
	assert.NotNil(t, accountAfter)

	assert.Equal(t, responseDto.Id, accountAfter.Id)
	assert.Equal(t, responseDto.Address, accountAfter.Address)
	assert.Equal(t, responseDto.Balance, accountAfter.Balance.String())
	assert.Equal(t, responseDto.Status, string(accountAfter.Status))
	assert.Equal(t, responseDto.CreatedAt, accountAfter.CreatedAt)
	assert.Equal(t, responseDto.UpdatedAt, accountAfter.UpdatedAt)
	assert.GreaterOrEqual(t, responseDto.CreatedAt, start)
	assert.GreaterOrEqual(t, responseDto.UpdatedAt, start)

	test.CompareAccount(t, accountAfter, responseDto)
}
