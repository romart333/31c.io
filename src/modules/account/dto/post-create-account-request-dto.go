package accountModuleDto

import (
	"fmt"
	errorHelpers "go-gin-test-job/src/common/error-helpers"
	errorMessages "go-gin-test-job/src/common/error-messages"
	"go-gin-test-job/src/common/validations"
	"go-gin-test-job/src/database/entities"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PostCreateAccountRequestDto struct {
	Address string                 `json:"address" validate:"AccountAddressValidation" example:"1JzfdUygUFk2M6KS3ngFMGRsy5vsH4N37a"`
	Name    string                 `json:"name" validate:"AccountNameValidation" example:"John Doe"`
	Rank    uint8                  `json:"rank" validate:"AccountRankValidation" example:"50"`
	Memo    string                 `json:"memo" example:"Some memo text"`
	Status  entities.AccountStatus `json:"status" validate:"AccountStatusValidation" enums:"On,Off" example:"On"`
}

var postCreateAccountRequestDtoValidator *validator.Validate

func init() {
	postCreateAccountRequestDtoValidator = validator.New()
	_ = postCreateAccountRequestDtoValidator.RegisterValidation("AccountAddressValidation", validations.AccountAddressValidation)
	_ = postCreateAccountRequestDtoValidator.RegisterValidation("AccountStatusValidation", validations.AccountStatusValidation)
	_ = postCreateAccountRequestDtoValidator.RegisterValidation("AccountRankValidation", validations.AccountRankValidation)
	_ = postCreateAccountRequestDtoValidator.RegisterValidation("AccountNameValidation", validations.AccountNameValidation)
}

func validatePostCreateAccountRequestDto(dto *PostCreateAccountRequestDto) error {
	return postCreateAccountRequestDtoValidator.Struct(dto)
}

// CreatePostCreateAccountRequestDto is the Gin version for handling the request
func CreatePostCreateAccountRequestDto(c *gin.Context) (PostCreateAccountRequestDto, error) {
	var dto PostCreateAccountRequestDto
	// Parse body params into DTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		errorMessage := PostCreateAccountRequestDtoQueryParseErrorMessage(err)
		return dto, errorHelpers.RespondBadRequestError(c, errorMessage)
	}
	// Validate the DTO
	if err := validatePostCreateAccountRequestDto(&dto); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errorMessage := PostCreateAccountRequestDtoValidateErrorMessage(err)
			return dto, errorHelpers.RespondBadRequestError(c, errorMessage)
		}
	}
	return dto, nil
}

func PostCreateAccountRequestDtoQueryParseErrorMessage(err error) string {
	return errorMessages.DefaultQueryParseErrorMessage()
}

func PostCreateAccountRequestDtoValidateErrorMessage(err validator.FieldError) string {
	var errorMessage string
	if err.Field() == "Address" && err.Tag() == "AccountAddressValidation" {
		errorMessage = fmt.Sprintf("%s format is wrong", err.Field())
	} else if err.Field() == "Status" && err.Tag() == "AccountStatusValidation" {
		errorMessage = fmt.Sprintf("%s must be one of the next values: %s", err.Field(), strings.Join(entities.AccountStatusList, ","))
	} else if err.Field() == "Rank" && err.Tag() == "AccountRankValidation" {
		errorMessage = fmt.Sprintf("%s must be between 0 and 100", err.Field())
	} else if err.Field() == "Name" && err.Tag() == "AccountNameValidation" {
		errorMessage = fmt.Sprintf("%s must be between 1 and 255 characters", err.Field())
	} else {
		errorMessage = errorMessages.DefaultFieldErrorMessage(err.Field())
	}
	return errorMessage
}
