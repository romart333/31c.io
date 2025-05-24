package accountTests

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	"github.com/stretchr/testify/assert"
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
	t.Run("TestGetAccountsRoute_SuccessParamsSearch", TestGetAccountsRoute_SuccessParamsSearch)
	t.Run("TestGetAccountsRoute_SuccessParamsSearchAndStatus", TestGetAccountsRoute_SuccessParamsSearchAndStatus)
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

	accounts, total := database.GetAccountsAndTotal("", make(map[string]string), accountModuleDto.DEFAULT_ACCOUNT_OFFSET, accountModuleDto.DEFAULT_ACCOUNT_COUNT, "")

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

	accounts, total := database.GetAccountsAndTotal("", make(map[string]string), params.Offset, params.Count, "")

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

	accounts, total := database.GetAccountsAndTotal(params.Status, make(map[string]string), accountModuleDto.DEFAULT_ACCOUNT_OFFSET, accountModuleDto.DEFAULT_ACCOUNT_COUNT, "")

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
	testCases := []struct {
		name    string
		orderBy string
	}{
		{
			"ByID",
			"id DESC",
		},
		{
			"ByUpdatedAt",
			"updated_at DESC",
		},
		{
			"ByAddress",
			"address ASC",
		},
		{
			"ByName",
			"name DESC",
		},
		{
			"ByRank",
			"rank ASC",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			type Params struct {
				OrderBy string `json:"orderBy"`
			}
			params := &Params{
				OrderBy: tc.orderBy,
			}

			query := url.Values{}
			query.Add("orderBy", params.OrderBy)

			u := &url.URL{
				Path:     fmt.Sprintf("/account"),
				RawQuery: query.Encode(),
			}

			orderParams, err := orderUtil.GetOrderByParamsSecure(nil, params.OrderBy, ",", accountModuleDto.GetAvailableAccountSortFieldList)
			accounts, total := database.GetAccountsAndTotal("", orderParams, accountModuleDto.DEFAULT_ACCOUNT_OFFSET, accountModuleDto.DEFAULT_ACCOUNT_COUNT, "")

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
		})
	}
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
	accounts, total := database.GetAccountsAndTotal(params.Status, orderParams, accountModuleDto.DEFAULT_ACCOUNT_OFFSET, accountModuleDto.DEFAULT_ACCOUNT_COUNT, "")

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
	accounts, total := database.GetAccountsAndTotal(params.Status, orderParams, params.Offset, params.Count, "")

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

func TestGetAccountsRoute_SuccessParamsSearch(t *testing.T) {
	type Params struct {
		Search string `json:"search"`
	}
	params := &Params{
		Search: "VIP",
	}

	query := url.Values{}
	query.Add("search", params.Search)

	u := &url.URL{
		Path:     fmt.Sprintf("/account"),
		RawQuery: query.Encode(),
	}

	accounts, total := database.GetAccountsAndTotal("", make(map[string]string), accountModuleDto.DEFAULT_ACCOUNT_OFFSET, accountModuleDto.DEFAULT_ACCOUNT_COUNT, params.Search)

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

	// Verify that all returned accounts contain the search term
	for _, accountDto := range responseDto.List {
		conditions := []func(account *entities.Account) bool{
			func(a *entities.Account) bool {
				return a.Id == accountDto.Id
			},
		}
		account := arrayUtil.FindItem(accounts, conditions)
		test.CompareAccount(t, *account, accountDto)

		// Check if the account matches the search criteria
		assert.True(t,
			strings.Contains(strings.ToLower((*account).Address), strings.ToLower(params.Search)) ||
				strings.Contains(strings.ToLower((*account).Name), strings.ToLower(params.Search)) ||
				strings.Contains(strings.ToLower((*account).Memo), strings.ToLower(params.Search)),
			"Account should match search criteria",
		)
	}
}

func TestGetAccountsRoute_SuccessParamsSearchAndStatus(t *testing.T) {
	type Params struct {
		Search string                 `json:"search"`
		Status entities.AccountStatus `json:"status"`
	}
	params := &Params{
		Search: "customer",
		Status: entities.AccountStatusOn,
	}

	query := url.Values{}
	query.Add("search", params.Search)
	query.Add("status", string(params.Status))

	u := &url.URL{
		Path:     fmt.Sprintf("/account"),
		RawQuery: query.Encode(),
	}

	accounts, total := database.GetAccountsAndTotal(params.Status, make(map[string]string), accountModuleDto.DEFAULT_ACCOUNT_OFFSET, accountModuleDto.DEFAULT_ACCOUNT_COUNT, params.Search)

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

	// Verify that all returned accounts contain the search term and match the status
	for _, accountDto := range responseDto.List {
		conditions := []func(account *entities.Account) bool{
			func(a *entities.Account) bool {
				return a.Id == accountDto.Id
			},
		}
		account := arrayUtil.FindItem(accounts, conditions)
		test.CompareAccount(t, *account, accountDto)

		// Check if the account matches the search criteria and status
		assert.True(t,
			strings.Contains(strings.ToLower((*account).Address), strings.ToLower(params.Search)) ||
				strings.Contains(strings.ToLower((*account).Name), strings.ToLower(params.Search)) ||
				strings.Contains(strings.ToLower((*account).Memo), strings.ToLower(params.Search)),
			"Account should match search criteria",
		)
		assert.Equal(t, params.Status, (*account).Status, "Account should match status filter")
	}
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
			accountModuleDto.PostCreateAccountRequestDto{
				Address: "invalid address",
				Name:    "John Doe",
				Rank:    50,
				Status:  entities.AccountStatusOn,
			},
			http.StatusBadRequest,
			errorHelpers.ResponseBadRequestErrorHTTP{Success: false, Message: "Address format is wrong"},
		},
		{
			"FailInvalidStatus",
			accountModuleDto.PostCreateAccountRequestDto{
				Address: "14yqg2y3a6HMgW9MiF5tVPAH4Dr1uxGKFJ",
				Name:    "John Doe",
				Rank:    50,
				Status:  "invalid status",
			},
			http.StatusBadRequest,
			errorHelpers.ResponseBadRequestErrorHTTP{Success: false, Message: fmt.Sprintf("%s must be one of the next values: %s", "Status", strings.Join(entities.AccountStatusList, ","))},
		},
		{
			"FailMissingName",
			accountModuleDto.PostCreateAccountRequestDto{
				Address: "14yqg2y3a6HMgW9MiF5tVPAH4Dr1uxGKFJ",
				Rank:    50,
				Status:  entities.AccountStatusOn,
			},
			http.StatusBadRequest,
			errorHelpers.ResponseBadRequestErrorHTTP{Success: false, Message: "Name must be between 1 and 255 characters"},
		},
		{
			"FailNameTooLong",
			accountModuleDto.PostCreateAccountRequestDto{
				Address: "14yqg2y3a6HMgW9MiF5tVPAH4Dr1uxGKFJ",
				Name:    strings.Repeat("a", 256),
				Rank:    50,
				Status:  entities.AccountStatusOn,
			},
			http.StatusBadRequest,
			errorHelpers.ResponseBadRequestErrorHTTP{Success: false, Message: "Name must be between 1 and 255 characters"},
		},
		{
			"FailMissingRank",
			accountModuleDto.PostCreateAccountRequestDto{
				Address: "14yqg2y3a6HMgW9MiF5tVPAH4Dr1uxGKFJ",
				Name:    "John Doe",
				Status:  entities.AccountStatusOn,
			},
			http.StatusBadRequest,
			errorHelpers.ResponseBadRequestErrorHTTP{Success: false, Message: "Rank must be between 0 and 100"},
		},
		{
			"FailRankTooHigh",
			accountModuleDto.PostCreateAccountRequestDto{
				Address: "14yqg2y3a6HMgW9MiF5tVPAH4Dr1uxGKFJ",
				Name:    "John Doe",
				Rank:    101,
				Status:  entities.AccountStatusOn,
			},
			http.StatusBadRequest,
			errorHelpers.ResponseBadRequestErrorHTTP{Success: false, Message: "Rank must be between 0 and 100"},
		},
	}

	for _, tt := range validationTests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.params)

			u := &url.URL{
				Path: fmt.Sprintf("/account"),
			}

			response := httptest.NewRecorder()
			request := httptest.NewRequest("POST", u.String(), bytes.NewBuffer(body))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("X-API-Key", config.AppConfig.AdminXApiKey)
			test.TestApp.ServeHTTP(response, request)
			assert.Equal(t, tt.expectedCode, response.Code)

			var responseBody errorHelpers.ResponseBadRequestErrorHTTP
			err := json.NewDecoder(response.Body).Decode(&responseBody)
			assert.Nil(t, err)
			assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
}

func TestCreateAccountRoute_FailAddressAlreadyExists(t *testing.T) {
	accountInfo := seeds.ACCOUNTS.ACCOUNT_1
	type Params struct {
		Address string                 `json:"address"`
		Name    string                 `json:"name"`
		Rank    uint8                  `json:"rank"`
		Memo    string                 `json:"memo"`
		Status  entities.AccountStatus `json:"status"`
	}
	params := &Params{
		Address: accountInfo.Address,
		Name:    "New Name",
		Rank:    60,
		Memo:    "New memo",
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

	// Verify that the existing account was not modified
	accountAfter := database.GetAccountByAddress(params.Address)
	assert.NotNil(t, accountAfter)
	assert.Equal(t, accountInfo.Name, accountAfter.Name, "Name should not be changed")
	assert.Equal(t, accountInfo.Rank, accountAfter.Rank, "Rank should not be changed")
	assert.Equal(t, accountInfo.Memo, accountAfter.Memo, "Memo should not be changed")
}

func TestCreateAccountRoute_Success(t *testing.T) {
	start := timeUtil.GetUnixTime()
	type Params struct {
		Address string                 `json:"address"`
		Name    string                 `json:"name"`
		Rank    uint8                  `json:"rank"`
		Memo    string                 `json:"memo"`
		Status  entities.AccountStatus `json:"status"`
	}
	params := &Params{
		Address: "32AaKxGbdhGMSGutcZjspFq9U89jJHW1um",
		Name:    "John Doe",
		Rank:    50,
		Memo:    "Test memo",
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
	assert.NotNil(t, responseDto.Name, "Name parameter should exist")
	assert.NotNil(t, responseDto.Rank, "Rank parameter should exist")
	assert.NotNil(t, responseDto.Memo, "Memo parameter should exist")
	assert.NotNil(t, responseDto.Balance, "Balance parameter should exist")
	assert.NotNil(t, responseDto.Status, "Status parameter should exist")
	assert.NotNil(t, responseDto.CreatedAt, "CreatedAt parameter should exist")
	assert.NotNil(t, responseDto.UpdatedAt, "UpdatedAt parameter should exist")

	accountAfter := database.GetAccountByAddress(params.Address)
	assert.NotNil(t, accountAfter)

	assert.Equal(t, responseDto.Id, accountAfter.Id)
	assert.Equal(t, responseDto.Address, accountAfter.Address)
	assert.Equal(t, responseDto.Name, accountAfter.Name)
	assert.Equal(t, responseDto.Rank, accountAfter.Rank)
	assert.Equal(t, responseDto.Memo, accountAfter.Memo)
	assert.Equal(t, responseDto.Balance, accountAfter.Balance.String())
	assert.Equal(t, responseDto.Status, string(accountAfter.Status))
	assert.Equal(t, responseDto.CreatedAt, accountAfter.CreatedAt)
	assert.Equal(t, responseDto.UpdatedAt, accountAfter.UpdatedAt)
	assert.GreaterOrEqual(t, responseDto.CreatedAt, start)
	assert.GreaterOrEqual(t, responseDto.UpdatedAt, start)

	test.CompareAccount(t, accountAfter, responseDto)
}
