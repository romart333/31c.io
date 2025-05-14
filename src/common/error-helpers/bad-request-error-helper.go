package errorHelpers

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type ResponseBadRequestErrorHTTP struct {
	Success bool   `json:"success" validate:"required" example:"false"`
	Message string `json:"message" validate:"required" example:"Bad request error"`
}

func NewResponseBadRequestErrorHTTP(message string) *ResponseBadRequestErrorHTTP {
	return &ResponseBadRequestErrorHTTP{
		Success: false,
		Message: message,
	}
}

func RespondBadRequestError(c *gin.Context, message string) error {
	if c != nil {
		c.JSON(400, NewResponseBadRequestErrorHTTP(message))
	}
	return fmt.Errorf("Bad request error. %s", message)
}
