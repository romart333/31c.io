package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-gin-test-job/src/database/entities"
	"go-gin-test-job/src/modules/account/dto"
	"go-gin-test-job/src/modules/account/service"
)

func TestGetAccount(t *testing.T) {
	// ... existing code ...

	t.Run("GetAccount_WithSorting", func(t *testing.T) {
		mockAccountService := new(service.MockAccountService)
		accountController := NewAccountController(mockAccountService)

		t.Run("GetAccount_WithSorting_ByAddress", func(t *testing.T) {
			request := &dto.GetAccountRequestDto{
				Page:  1,
				Limit: 10,
				Sort:  "address",
			}

			expectedAccounts := []entities.Account{
				{Id: 1, Address: "0x123", Name: "Test1", Rank: 1},
				{Id: 2, Address: "0x456", Name: "Test2", Rank: 2},
			}

			mockAccountService.On("GetAccounts", mock.Anything, request).Return(expectedAccounts, int64(2), nil)

			accounts, total, err := accountController.GetAccount(request)

			assert.NoError(t, err)
			assert.Equal(t, int64(2), total)
			assert.Equal(t, expectedAccounts, accounts)
			mockAccountService.AssertExpectations(t)
		})

		t.Run("GetAccount_WithSorting_ByName", func(t *testing.T) {
			request := &dto.GetAccountRequestDto{
				Page:  1,
				Limit: 10,
				Sort:  "name",
			}

			expectedAccounts := []entities.Account{
				{Id: 1, Address: "0x123", Name: "Test1", Rank: 1},
				{Id: 2, Address: "0x456", Name: "Test2", Rank: 2},
			}

			mockAccountService.On("GetAccounts", mock.Anything, request).Return(expectedAccounts, int64(2), nil)

			accounts, total, err := accountController.GetAccount(request)

			assert.NoError(t, err)
			assert.Equal(t, int64(2), total)
			assert.Equal(t, expectedAccounts, accounts)
			mockAccountService.AssertExpectations(t)
		})

		t.Run("GetAccount_WithSorting_ByRank", func(t *testing.T) {
			request := &dto.GetAccountRequestDto{
				Page:  1,
				Limit: 10,
				Sort:  "rank",
			}

			expectedAccounts := []entities.Account{
				{Id: 1, Address: "0x123", Name: "Test1", Rank: 1},
				{Id: 2, Address: "0x456", Name: "Test2", Rank: 2},
			}

			mockAccountService.On("GetAccounts", mock.Anything, request).Return(expectedAccounts, int64(2), nil)

			accounts, total, err := accountController.GetAccount(request)

			assert.NoError(t, err)
			assert.Equal(t, int64(2), total)
			assert.Equal(t, expectedAccounts, accounts)
			mockAccountService.AssertExpectations(t)
		})

		t.Run("GetAccount_WithSorting_ByInvalidField", func(t *testing.T) {
			request := &dto.GetAccountRequestDto{
				Page:  1,
				Limit: 10,
				Sort:  "invalid_field",
			}

			expectedAccounts := []entities.Account{
				{Id: 1, Address: "0x123", Name: "Test1", Rank: 1},
				{Id: 2, Address: "0x456", Name: "Test2", Rank: 2},
			}

			mockAccountService.On("GetAccounts", mock.Anything, request).Return(expectedAccounts, int64(2), nil)

			accounts, total, err := accountController.GetAccount(request)

			assert.NoError(t, err)
			assert.Equal(t, int64(2), total)
			assert.Equal(t, expectedAccounts, accounts)
			mockAccountService.AssertExpectations(t)
		})
	})

	// ... existing code ...
}
