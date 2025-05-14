package blockchain

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"go-gin-test-job/src/config"
	currencyUtil "go-gin-test-job/src/utils/currency"
	timeUtil "go-gin-test-job/src/utils/time"
	"net/http"
)

var externalUrl = "https://api.bitcore.io/api/BTC/mainnet"

type BlockchainBalanceResponse struct {
	Confirmed int64 `json:"confirmed"`
}

func GetAddressBalance(address string) (decimal.Decimal, error) {
	balance := decimal.NewFromInt(0)
	url := fmt.Sprintf("%s/address/%s/balance", externalUrl, address)
	client := &http.Client{
		Timeout: timeUtil.DurationSeconds(config.AppConfig.RequestTimeoutSec),
	}
	response, err := client.Get(url)
	if err != nil {
		return balance, err
	}
	defer response.Body.Close()
	var responseData BlockchainBalanceResponse
	if err := json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		return balance, err
	}
	balance = currencyUtil.FromSatoshi(responseData.Confirmed)
	return balance, nil
}
