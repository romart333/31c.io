package validations

import (
	"go-gin-test-job/src/database/entities"
	addressValidationUtil "go-gin-test-job/src/utils/address-validation"
	nameValidationUtil "go-gin-test-job/src/utils/name-validation"
	rankValidationUtil "go-gin-test-job/src/utils/rank-validation"
	"strings"

	"github.com/go-playground/validator/v10"
)

func AccountStatusValidation(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	status = strings.Trim(status, "\"")
	switch entities.AccountStatus(status) {
	case entities.AccountStatusOn, entities.AccountStatusOff:
		return true
	}
	return false
}

func AccountAddressValidation(fl validator.FieldLevel) bool {
	address := fl.Field().String()
	return addressValidationUtil.IsValidAddress(address)
}

func NotEmpty(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	return strings.TrimSpace(str) != ""
}

func AccountRankValidation(fl validator.FieldLevel) bool {
	rank := fl.Field().Uint()
	return rankValidationUtil.IsValidRank(rank)
}

func AccountNameValidation(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	return nameValidationUtil.IsValidName(name)
}
