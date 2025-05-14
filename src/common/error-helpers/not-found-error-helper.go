package errorHelpers

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type ResponseNotFoundErrorHTTP struct {
	Success bool   `json:"success" validate:"required" example:"false"`
	Message string `json:"message" validate:"required" example:"Not found error"`
}

func NewResponseNotFoundErrorHTTP(message string) *ResponseNotFoundErrorHTTP {
	return &ResponseNotFoundErrorHTTP{
		Success: false,
		Message: message,
	}
}

func RespondNotFoundError(c *gin.Context, message string) error {
	if c != nil {
		c.JSON(404, NewResponseNotFoundErrorHTTP(message))
	}
	return fmt.Errorf("Not found error. %s", message)
}
