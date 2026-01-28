package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jekiti/citydrive/api-gateway/internal/common"
)

func RequireCarAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := common.LoggerForModule(c, "middleware", "RequireCarAuth")

		traceID := common.GetTraceID(c)
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			logger.Warn("missing authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error_code":        "MISSING_AUTH_HEADER",
				"error_description": "Authorization header is required",
				"trace_id":          traceID,
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			logger.Warn("invalid authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error_code":        "INVALID_AUTH_HEADER",
				"error_description": "Authorization header format must be Bearer {token}",
				"trace_id":          traceID,
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			logger.Warn("invalid token", "error", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error_code":        "INVALID_TOKEN",
				"error_description": "Invalid token",
				"trace_id":          traceID,
			})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				logger.Warn("token expired")
				c.JSON(http.StatusUnauthorized, gin.H{
					"error_code":        "TOKEN_EXPIRED",
					"error_description": "Token is expired",
					"trace_id":          traceID,
				})
				c.Abort()
				return
			}
			roles, ok := claims["roles"].([]any)
			if !ok {
				logger.Warn("invalid token roles claim")
				c.JSON(http.StatusUnauthorized, gin.H{
					"error_code":        "INVALID_CLAIMS",
					"error_description": "Invalid token claims",
					"trace_id":          traceID,
				})
				c.Abort()
				return
			}
			if roles[0] != "car" {
				logger.Warn("insufficient token roles", "roles", roles)
				c.JSON(http.StatusForbidden, gin.H{
					"error_code":        "INSUFFICIENT_ROLES",
					"error_description": "Insufficient token roles",
					"trace_id":          traceID,
				})
				c.Abort()
				return
			}

			c.Set("car_id", claims["car_id"])

			c.Set("claims", claims)

			logger.Info("car authenticated", "car_id", claims["car_id"])
			c.Next()
		} else {
			logger.Warn("invalid token claims")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error_code":        "INVALID_CLAIMS",
				"error_description": "Invalid token claims",
				"trace_id":          traceID,
			})
			c.Abort()
			return
		}
	}
}
