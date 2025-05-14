package addressValidationUtil

import (
	"regexp"
)

var addressRegex = regexp.MustCompile(`^[13][a-km-zA-HJ-NP-Z0-9]{26,33}$`)

func IsValidAddress(address string) bool {
	return addressRegex.MatchString(address)
}
