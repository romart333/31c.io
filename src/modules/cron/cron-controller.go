package cronModule

import (
	"github.com/gin-gonic/gin"
	"go-gin-test-job/src/common/dto"
)

// UpdateAccountsBalances Update accounts balances
// @Summary Update accounts balances
// @Description Update accounts balances
// @Tags Cron
// @Accept json
// @Produce json
// @Param X-API-Key header string true "Cron api key"
// @Success 201 {object} dto.SuccessDto
// @Failure 401 {object} errorHelpers.ResponseUnauthorizedErrorHTTP{}
// @Router /cron/account-balance [post]
func UpdateAccountsBalances(c *gin.Context) {
	updateAccountsBalances()
	c.JSON(200, dto.CreateSuccessDto())
}
