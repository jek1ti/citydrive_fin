package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jekiti/citydrive/api-gateway/internal/common"
)

func RequireAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := common.LoggerForModule(c, "middleware", "RequireAuth")

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
			exp, ok := claims["exp"].(float64)
			if !ok || float64(time.Now().Unix()) > exp {
				logger.Warn("token expired")
				c.JSON(http.StatusUnauthorized, gin.H{
					"error_code":        "TOKEN_EXPIRED",
					"error_description": "Token is expired",
					"trace_id":          traceID,
				})
				c.Abort()
				return
			}

			c.Set("user_id", claims["sub"])
			c.Set("email", claims["email"])
			c.Set("roles", claims["roles"])
			c.Set("claims", claims)

			logger.Info("user authenticated", "user_id", claims["sub"])
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
