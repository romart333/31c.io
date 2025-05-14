package accountModuleDto

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	errorHelpers "go-gin-test-job/src/common/error-helpers"
	errorMessages "go-gin-test-job/src/common/error-messages"
	"go-gin-test-job/src/common/validations"
	"go-gin-test-job/src/database/entities"
	stringUtil "go-gin-test-job/src/utils/string"
	"strings"
)

const DEFAULT_ACCOUNT_COUNT = 100
const DEFAULT_ACCOUNT_OFFSET = 0

var GetAvailableAccountSortField = map[string]string{
	"id":         "account.id",
	"updated_at": "account.updated_at",
}

var GetAvailableAccountSortFieldList = func() []string {
	keys := make([]string, 0, len(GetAvailableAccountSortField))
	for key := range GetAvailableAccountSortField {
		keys = append(keys, key)
	}
	return keys
}()

type GetAccountRequestDto struct {
	Offset  int                    `form:"offset" json:"offset" validate:"min=0" default:"0" example:"5"`
	Count   int                    `form:"count" json:"count" validate:"min=1,max=100" default:"100" example:"20"`
	Status  entities.AccountStatus `form:"status" json:"status" validate:"omitempty,AccountStatusValidation" example:"On"`
	OrderBy string                 `form:"orderBy" json:"orderBy" validate:"omitempty,max=255" example:"id ASC"`
}

var getAccountRequestDtoValidator *validator.Validate

func init() {
	getAccountRequestDtoValidator = validator.New()
	_ = getAccountRequestDtoValidator.RegisterValidation("AccountStatusValidation", validations.AccountStatusValidation)
}

func getAccountRequestDtoDefaultValues(dto *GetAccountRequestDto) {
	if dto.Count == 0 {
		dto.Count = DEFAULT_ACCOUNT_COUNT
	}
}

func validateGetAccountRequestDto(dto *GetAccountRequestDto) error {
	return getAccountRequestDtoValidator.Struct(dto)
}

// CreateGetAccountRequestDto is the Gin version of handling the request
func CreateGetAccountRequestDto(c *gin.Context) (GetAccountRequestDto, error) {
	var dto GetAccountRequestDto
	// Parse query params into DTO
	if err := c.ShouldBindQuery(&dto); err != nil {
		errorMessage := GetAccountRequestDtoQueryParseErrorMessage(err)
		return dto, errorHelpers.RespondBadRequestError(c, errorMessage)
	}
	// Set default values
	getAccountRequestDtoDefaultValues(&dto)
	// Validate the DTO
	if err := validateGetAccountRequestDto(&dto); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errorMessage := GetAccountRequestDtoValidateErrorMessage(err)
			return dto, errorHelpers.RespondBadRequestError(c, errorMessage)
		}
	}
	dto.Status = entities.AccountStatus(strings.Trim(string(dto.Status), "\""))
	return dto, nil
}

func GetAccountRequestDtoQueryParseErrorMessage(err error) string {
	var errorMessage string
	if stringUtil.CaseInsensitiveContains(err.Error(), "\"offset\"") || stringUtil.CaseInsensitiveContains(err.Error(), ".offset") {
		errorMessage = errorMessages.DefaultFieldErrorMessage("offset")
	} else if stringUtil.CaseInsensitiveContains(err.Error(), "\"count\"") || stringUtil.CaseInsensitiveContains(err.Error(), ".count") {
		errorMessage = errorMessages.DefaultFieldErrorMessage("count")
	} else {
		errorMessage = errorMessages.DefaultQueryParseErrorMessage()
	}
	return errorMessage
}

func GetAccountRequestDtoValidateErrorMessage(err validator.FieldError) string {
	var errorMessage string
	if err.Field() == "Count" && err.Tag() == "min" {
		errorMessage = fmt.Sprintf("%s must be greater than or equal %s", err.Field(), err.Param())
	} else if err.Field() == "Count" && err.Tag() == "max" {
		errorMessage = fmt.Sprintf("%s must be less than or equal %s", err.Field(), err.Param())
	} else if err.Field() == "Offset" && err.Tag() == "min" {
		errorMessage = fmt.Sprintf("%s must be greater than or equal %s", err.Field(), err.Param())
	} else if err.Field() == "Status" && err.Tag() == "AccountStatusValidation" {
		errorMessage = fmt.Sprintf("%s must be one of the next values: %s", err.Field(), strings.Join(entities.AccountStatusList, ","))
	} else if err.Field() == "OrderBy" && err.Tag() == "max" {
		errorMessage = fmt.Sprintf("%s must be shorter than or equal to %s characters", err.Field(), err.Param())
	} else {
		errorMessage = errorMessages.DefaultFieldErrorMessage(err.Field())
	}
	return errorMessage
}
