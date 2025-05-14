package errorMessages

import (
	"fmt"
)

func DefaultFieldErrorMessage(field string) string {
	return fmt.Sprintf("%s is invalid", field)
}

func DefaultQueryParseErrorMessage() string {
	return "Invalid request query"
}
