package common

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func Response(c *gin.Context, httpCode int, code, description, err string) {
	traceID := GetTraceID(c)
	c.JSON(httpCode, gin.H{
		"error_code":        code,
		"error_description": description,
		"trace_id":          traceID,
		"error":             err,
	})
}

func LoggerForModule(c *gin.Context, module, function string) *slog.Logger {
	traceID := GetTraceID(c)

	return slog.Default().With(
		"module", module,
		"func", function,
		"trace_id", traceID,
		"time", time.Now().Format(time.RFC3339),
	)
}

func LoggerWithID(c *gin.Context, module, function, id string) *slog.Logger {
	logger := LoggerForModule(c, module, function)
	return logger.With("id", id)
}

func GetTraceID(c *gin.Context) string {
	if traceID, exists := c.Get("trace_id"); exists {
		return traceID.(string)
	}
	return ""
}
