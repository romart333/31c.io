package accountModule

import (
	"github.com/gin-gonic/gin"
	accountModuleDto "go-gin-test-job/src/modules/account/dto"
	orderUtil "go-gin-test-job/src/utils/order"
)

// GetAccounts Get list of accounts
// @Summary Get list of accounts
// @Description Get list of account
// @Tags Account
// @Accept json
// @Produce json
// @Param offset query int false "This is paging offset. 0 by default" minimum(0) default(0)
// @Param count query int false "Max item count in single response. 100 by default" minimum(1) maximum(100) default(100)
// @Param status query string false "Account statuses: On, Off" Enums("On", "Off") default("On")
// @Param orderBy query string false "Comma-separated sort order options (sort fields: id, updated, sort order: ASC,DESC)" default(id ASC)
// @Param X-API-Key header string true "Admin api key"
// @Success 200 {object} accountModuleDto.GetAccountResponseDto
// @Failure 400 {object} errorHelpers.ResponseBadRequestErrorHTTP{}
// @Failure 401 {object} errorHelpers.ResponseUnauthorizedErrorHTTP{}
// @Router /account [get]
func GetAccounts(c *gin.Context) {
	dto, err := accountModuleDto.CreateGetAccountRequestDto(c)
	if err != nil {
		return
	}
	orderParams, err := orderUtil.GetOrderByParamsSecure(c, dto.OrderBy, ",", accountModuleDto.GetAvailableAccountSortFieldList)
	if err != nil {
		return
	}
	accounts, total := getAccounts(dto.Status, orderParams, dto.Offset, dto.Count)
	c.JSON(200, accountModuleDto.CreateGetAccountResponseDto(dto.Offset, dto.Count, total, accounts))
}

// CreateAccount Create new account
// @Summary Create new account
// @Description Create new account
// @Tags Account
// @Accept json
// @Produce json
// @Param X-API-Key header string true "Admin api key"
// @Param request body accountModuleDto.PostCreateAccountRequestDto true "Request body"
// @Success 200 {object} accountModuleDto.AccountDto
// @Failure 400 {object} errorHelpers.ResponseBadRequestErrorHTTP{}
// @Failure 401 {object} errorHelpers.ResponseUnauthorizedErrorHTTP{}
// @Failure 409 {object} errorHelpers.ResponseConflictErrorHTTP{}
// @Router /account [post]
func CreateAccount(c *gin.Context) {
	dto, err := accountModuleDto.CreatePostCreateAccountRequestDto(c)
	if err != nil {
		return
	}
	account, err := createAccount(c, dto.Address, dto.Status)
	if err != nil {
		return
	}
	c.JSON(200, accountModuleDto.CreatePostCreateAccountResponseDto(account))
}
