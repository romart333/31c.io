package validations

import (
	"github.com/go-playground/validator/v10"
	"go-gin-test-job/src/database/entities"
	addressValidationUtil "go-gin-test-job/src/utils/address-validation"
	"strings"
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
	if addressValidationUtil.IsValidAddress(address) {
		return true
	}
	return false
}

func NotEmpty(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	return strings.TrimSpace(str) != ""
}
