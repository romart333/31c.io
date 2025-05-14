package cronModule

import (
	"fmt"
	"go-gin-test-job/src/config"
	"go-gin-test-job/src/database"
	"go-gin-test-job/src/database/entities"
	"go-gin-test-job/src/logger"
	"go-gin-test-job/src/modules/common/blockchain"
)

func updateAccountsBalances() {
	accounts := database.GetAccountsBatch(config.AppConfig.CronBatchCount)
	for _, account := range accounts {
		if err := updateAccountBalance(account); err != nil {
			logger.Logger.Error().Msg(fmt.Sprintf("Update account %d address %s error. %s", account.Id, account.Address, err.Error()))
		}
	}
}

func updateAccountBalance(account *entities.Account) error {
	logger.Logger.Info().Msg(fmt.Sprintf("Update account %d address %s balance", account.Id, account.Address))
	balance, err := blockchain.GetAddressBalance(account.Address)
	if err != nil {
		return err
	}
	logger.Logger.Info().Msg(fmt.Sprintf("Account %d address %s balance - %s", account.Id, account.Address, account.Balance))
	updateData := account.UpdateBalance(balance)
	if err := database.UpdateAccount(nil, account, updateData); err != nil {
		return err
	}
	return nil
}
