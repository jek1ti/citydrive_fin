package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GenerateRID(id string, timestamp time.Time) string {
	utcTimestamp := timestamp.UTC().Unix()

	data := fmt.Sprintf("%s%s",
		id,
		strconv.FormatInt(utcTimestamp, 10))

	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:16])
}

func TracingMiddleware(headerName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := GenerateRID("api-gateway", time.Now())

		c.Set("trace_id", traceID)

		c.Header(headerName, traceID)
		
		c.Next()
	}
}

