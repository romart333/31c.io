package errorHelpers

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type ResponseUnauthorizedErrorHTTP struct {
	Success bool   `json:"success" validate:"required" example:"false"`
	Message string `json:"message" validate:"required" example:"Unauthorized error"`
}

func NewResponseUnauthorizedErrorHTTP(message string) *ResponseUnauthorizedErrorHTTP {
	return &ResponseUnauthorizedErrorHTTP{
		Success: false,
		Message: message,
	}
}

func RespondUnauthorizedError(c *gin.Context) error {
	if c != nil {
		c.JSON(401, NewResponseUnauthorizedErrorHTTP("Unauthorized"))
	}
	return fmt.Errorf("Unauthorized error")
}
