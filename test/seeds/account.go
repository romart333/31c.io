package seeds

import (
	"github.com/shopspring/decimal"
	"go-gin-test-job/src/database/entities"
	timeUtil "go-gin-test-job/src/utils/time"
)

var ACCOUNTS struct {
	ACCOUNT_1 entities.Account
	ACCOUNT_2 entities.Account
	ACCOUNT_3 entities.Account
	ACCOUNT_4 entities.Account
}

func FillAccountList() []entities.Account {
	ACCOUNTS.ACCOUNT_1 = entities.Account{
		Id:        1,
		Address:   "3JTCWLKubxuuXXnmQPxx43nP2LJAcPSL1W",
		Balance:   decimal.RequireFromString("0.96224397"),
		Status:    entities.AccountStatusOn,
		CreatedAt: timeUtil.GetUnixTime(),
		UpdatedAt: timeUtil.GetUnixTime(),
	}
	ACCOUNTS.ACCOUNT_2 = entities.Account{
		Id:        2,
		Address:   "38JeTiYSS2Y4kSxNBNH6kmH5kjm8sodDvU",
		Balance:   decimal.RequireFromString("0.00056665"),
		Status:    entities.AccountStatusOn,
		CreatedAt: timeUtil.GetUnixTime(),
		UpdatedAt: timeUtil.GetUnixTime(),
	}
	ACCOUNTS.ACCOUNT_3 = entities.Account{
		Id:        3,
		Address:   "34bMmbjiiK5WfV2ZtgZGxLVYycJGNPEqjE",
		Balance:   decimal.NewFromInt(0),
		Status:    entities.AccountStatusOff,
		CreatedAt: timeUtil.GetUnixTime(),
		UpdatedAt: timeUtil.GetUnixTime(),
	}
	ACCOUNTS.ACCOUNT_4 = entities.Account{
		Id:        4,
		Address:   "1CmSPVJifmK3HXqy2tYgbTSb4eExK4wqYT",
		Balance:   decimal.RequireFromString("0.07134313"),
		Status:    entities.AccountStatusOff,
		CreatedAt: timeUtil.GetUnixTime(),
		UpdatedAt: timeUtil.GetUnixTime(),
	}
	return []entities.Account{
		ACCOUNTS.ACCOUNT_1,
		ACCOUNTS.ACCOUNT_2,
		ACCOUNTS.ACCOUNT_3,
		ACCOUNTS.ACCOUNT_4,
	}
}

func GetAccountList() []entities.Account {
	return []entities.Account{
		ACCOUNTS.ACCOUNT_1,
		ACCOUNTS.ACCOUNT_2,
		ACCOUNTS.ACCOUNT_3,
		ACCOUNTS.ACCOUNT_4,
	}
}
