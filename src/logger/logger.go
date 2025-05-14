package logger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"time"
)

var Logger zerolog.Logger

func InitializeLogger() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	consoleWriter := zerolog.NewConsoleWriter()
	consoleWriter.FormatLevel = func(i interface{}) string {
		switch i {
		case "info":
			return fmt.Sprintf("\033[32m%s\033[0m", i) // Green for info
		case "warn":
			return fmt.Sprintf("\033[33m%s\033[0m", i) // Yellow for warn
		case "error":
			return fmt.Sprintf("\033[31m%s\033[0m", i) // Red for error
		case "fatal":
			return fmt.Sprintf("\033[41m%s\033[0m", i) // White text on red background for fatal
		default:
			return fmt.Sprintf("\033[37m%s\033[0m", i) // Default color for other levels
		}
	}
	Logger = zerolog.New(consoleWriter).With().Timestamp().Logger()
}

func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := c.GetHeader("X-Request-ID")
		c.Next()
		event := Logger.Info().
			Str("requestid", requestID).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Dur("duration", time.Since(start))
		if len(c.Errors) > 0 {
			event.Msg("Request failed. Error - " + c.Errors.String())
		} else {
			event.Msg("Request completed")
		}
	}
}

func SetDebugLevel() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}
