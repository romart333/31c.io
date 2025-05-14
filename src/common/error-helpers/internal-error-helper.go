package errorHelpers

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type ResponseInternalErrorHTTP struct {
	Success bool   `json:"success" validate:"required" example:"false"`
	Message string `json:"message" validate:"required" example:"Internal error"`
}

func NewResponseInternalErrorHTTP(message string) *ResponseInternalErrorHTTP {
	return &ResponseInternalErrorHTTP{
		Success: false,
		Message: message,
	}
}

func RespondInternalError(c *gin.Context, message string) error {
	if c != nil {
		c.JSON(500, NewResponseInternalErrorHTTP(message))
	}
	return fmt.Errorf("Internal error. %s", message)
}
