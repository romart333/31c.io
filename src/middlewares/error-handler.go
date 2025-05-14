package middleware

import (
	"github.com/gin-gonic/gin"
	errorHelpers "go-gin-test-job/src/common/error-helpers"
	"net/http"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			if err != nil {
				switch c.Writer.Status() {
				case http.StatusBadRequest:
					c.JSON(http.StatusBadRequest, errorHelpers.NewResponseBadRequestErrorHTTP(err.Error()))
				case http.StatusUnauthorized:
					c.JSON(http.StatusUnauthorized, errorHelpers.NewResponseUnauthorizedErrorHTTP(err.Error()))
				case http.StatusNotFound:
					c.JSON(http.StatusNotFound, errorHelpers.NewResponseNotFoundErrorHTTP(err.Error()))
				case http.StatusConflict:
					c.JSON(http.StatusConflict, errorHelpers.NewResponseConflictErrorHTTP(err.Error()))
				case http.StatusInternalServerError:
					c.JSON(http.StatusInternalServerError, errorHelpers.NewResponseInternalErrorHTTP(err.Error()))
				default:
					c.JSON(http.StatusInternalServerError, errorHelpers.NewResponseInternalErrorHTTP(err.Error()))
				}
			}
		}
	}
}
