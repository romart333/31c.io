package orderUtil

import (
	"fmt"
	"github.com/gin-gonic/gin"
	errorHelpers "go-gin-test-job/src/common/error-helpers"
	"strings"
)

var AvailableSortOrderList = map[string]bool{
	"ASC":  true,
	"DESC": true,
}

func GetOrderByParamsSecure(c *gin.Context, data, separator string, availableSortFieldList []string) (map[string]string, error) {
	orderByResult := make(map[string]string)
	availableSortFields := make(map[string]bool)
	// Convert availableSortFieldList to a map for faster lookup
	for _, field := range availableSortFieldList {
		availableSortFields[field] = true
	}
	// Split and process the order-by parameters
	orderByList := strings.Split(data, separator)
	for _, orderByLine := range orderByList {
		orderByLine = strings.TrimSpace(orderByLine)
		if orderByLine == "" {
			continue
		}
		parts := strings.Fields(orderByLine) // Splits by whitespace
		if len(parts) < 1 || len(parts) > 2 {
			// Return a structured bad request error
			return nil, errorHelpers.RespondBadRequestError(c, fmt.Sprintf("invalid order by parameter: %s", orderByLine))
		}
		order := parts[0]
		direction := "ASC" // Default sort order
		if len(parts) == 2 {
			direction = strings.ToUpper(parts[1])
		}
		// Validate order field
		if !availableSortFields[order] {
			// Return a structured bad request error

			return nil, errorHelpers.RespondBadRequestError(c, fmt.Sprintf("cannot order by %s", orderByLine))
		}
		// Validate direction
		if !AvailableSortOrderList[direction] {
			// Return a structured bad request error
			return nil, errorHelpers.RespondBadRequestError(c, fmt.Sprintf("invalid order direction: %s", direction))
		}
		// Avoid duplicate order fields
		if _, exists := orderByResult[order]; !exists {
			orderByResult[order] = direction
		}
	}

	return orderByResult, nil
}
