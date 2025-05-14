package errorHelpers

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type ResponseConflictErrorHTTP struct {
	Success bool   `json:"success" validate:"required" example:"false"`
	Message string `json:"message" validate:"required" example:"Conflict error"`
}

func NewResponseConflictErrorHTTP(message string) *ResponseConflictErrorHTTP {
	return &ResponseConflictErrorHTTP{
		Success: false,
		Message: message,
	}
}

func RespondConflictError(c *gin.Context, message string) error {
	if c != nil {
		c.JSON(409, NewResponseConflictErrorHTTP(message))
	}
	return fmt.Errorf("Conflict error. %s", message)
}
