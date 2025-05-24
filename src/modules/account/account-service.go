package accountModule

import (
	errorHelpers "go-gin-test-job/src/common/error-helpers"
	"go-gin-test-job/src/database"
	"go-gin-test-job/src/database/entities"
	accountModuleDto "go-gin-test-job/src/modules/account/dto"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func getAccounts(status entities.AccountStatus, orderParams map[string]string, offset int, count int, search string) ([]*entities.Account, int64) {
	return database.GetAccountsAndTotal(status, orderParams, offset, count, search)
}

func createAccount(c *gin.Context, dto accountModuleDto.PostCreateAccountRequestDto) (*entities.Account, error) {
	var account *entities.Account
	transactionError := database.DbConn.Transaction(func(tx *gorm.DB) error {
		if database.IsAddressExists(tx, dto.Address) {
			return errorHelpers.RespondConflictError(c, "Address already exists")
		}
		newAccount := entities.CreateAccount(dto.Address, dto.Status, dto.Name, dto.Rank, dto.Memo)
		var err error
		account, err = database.CreateAccount(tx, newAccount)
		if err != nil {
			return err
		}
		return nil
	}, database.DefaultTxOptions)
	if transactionError != nil {
		return nil, transactionError
	}
	return account, nil
}
