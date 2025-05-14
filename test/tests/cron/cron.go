package cronTests

import (
	"encoding/json"
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go-gin-test-job/src/common/dto"
	"go-gin-test-job/src/config"
	"go-gin-test-job/src/database"
	"go-gin-test-job/src/database/entities"
	arrayUtil "go-gin-test-job/src/utils/array"
	currencyUtil "go-gin-test-job/src/utils/currency"
	numberUtil "go-gin-test-job/src/utils/number"
	timeUtil "go-gin-test-job/src/utils/time"
	"go-gin-test-job/test"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestCronRoute(t *testing.T) {
	t.Run("TestUpdateAccountsBalancesRoute_Success", TestUpdateAccountsBalancesRoute_Success)
}

func TestUpdateAccountsBalancesRoute_Success(t *testing.T) {
	start := timeUtil.GetUnixTime()

	u := &url.URL{
		Path: fmt.Sprintf("/cron/account-balance"),
	}

	accountsBefore := database.GetAccountsBatch(config.AppConfig.CronBatchCount)
	assert.Greater(t, len(accountsBefore), 0)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockAccountsBalance := make(map[int64]decimal.Decimal)
	for _, accountBefore := range accountsBefore {
		mockBalance := int64(numberUtil.GetRandomNumber(0, 10000000000))
		// Define the mock response
		httpmock.RegisterResponder(
			"GET",
			fmt.Sprintf("https://api.bitcore.io/api/BTC/mainnet/address/%s/balance", accountBefore.Address),
			httpmock.NewStringResponder(200, fmt.Sprintf(`{"confirmed": %d}`, mockBalance)),
		)
		mockAccountsBalance[accountBefore.Id] = currencyUtil.FromSatoshi(mockBalance)
	}

	response := httptest.NewRecorder()
	request := httptest.NewRequest("POST", u.String(), nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-Key", config.AppConfig.CronXApiKey)
	test.TestApp.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code)

	// Read the response body and parse JSON
	var responseDto dto.SuccessDto
	err := json.NewDecoder(response.Body).Decode(&responseDto)
	assert.Nil(t, err)

	assert.NotNil(t, responseDto.Success, "Success parameter should exist")

	assert.Equal(t, true, responseDto.Success)

	accountIds := make([]int64, 0)
	for _, account := range accountsBefore {
		accountIds = append(accountIds, account.Id)
	}

	accountsAfter := database.GetAccountsByIds(accountIds)
	assert.Equal(t, len(accountsBefore), len(accountsAfter))

	for _, accountAfter := range accountsAfter {
		conditions := []func(account *entities.Account) bool{
			func(a *entities.Account) bool {
				return a.Id == accountAfter.Id
			},
		}
		accountBefore := arrayUtil.FindItem(accountsBefore, conditions)
		assert.NotNil(t, accountBefore)

		assert.Equal(t, (*accountBefore).Id, accountAfter.Id)
		assert.Equal(t, (*accountBefore).Address, accountAfter.Address)
		assert.Equal(t, mockAccountsBalance[accountAfter.Id].String(), accountAfter.Balance.String())
		assert.Equal(t, (*accountBefore).CreatedAt, accountAfter.CreatedAt)
		assert.GreaterOrEqual(t, accountAfter.UpdatedAt, (*accountBefore).UpdatedAt)
		assert.GreaterOrEqual(t, accountAfter.UpdatedAt, start)
	}
}
