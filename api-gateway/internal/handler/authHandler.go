package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jekiti/citydrive/api-gateway/internal/common"
	"github.com/jekiti/citydrive/api-gateway/internal/model"
	"github.com/jekiti/citydrive/api-gateway/internal/service"
	authpb "github.com/jekiti/citydrive/gen/proto/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	authClient *service.AuthClient
}

func NewAuthHandler(authClient *service.AuthClient) *AuthHandler {
	return &AuthHandler{authClient: authClient}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Response(c, 400, "INVALID_REQUEST", "Invalid request payload", err.Error())
		return
	}
	traceID := common.GetTraceID(c)

	ctx := c.Request.Context()
	resp, err := h.authClient.Login(ctx, traceID, &authpb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.Unavailable:
			common.Response(c, 502, "SERVICE_UNAVAILABLE", "Auth service is down", err.Error())
			return
		case codes.DeadlineExceeded:
			common.Response(c, 504, "TIMEOUT", "Request timeout", err.Error())
			return
		case codes.InvalidArgument:
			common.Response(c, 400, "INVALID_DATA", "Invalid Auth data", err.Error())
			return
		case codes.PermissionDenied:
			common.Response(c, 403, "PERMISSION_DENIED", "Access denied", err.Error())
			return
		default:
			common.Response(c, 500, "INTERNAL_ERROR", "Internal server error", err.Error())
			return
		}
	}
	c.JSON(200, gin.H{
		"access_token": resp.AccessToken,
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Response(c, 400, "invalid_request", "Invalid request payload", err.Error())
		return
	}

	traceID := common.GetTraceID(c)
	ctx := c.Request.Context()
	resp, err := h.authClient.Register(ctx, traceID, &authpb.RegisterRequest{
		Email:      req.Email,
		Name:       req.Name,
		Surname:    req.Surname,
		Department: req.Department,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.Unavailable:
			common.Response(c, 502, "SERVICE_UNAVAILABLE", "Auth service is down", err.Error())
			return
		case codes.DeadlineExceeded:
			common.Response(c, 504, "TIMEOUT", "Request timeout", err.Error())
			return
		case codes.InvalidArgument:
			common.Response(c, 400, "INVALID_DATA", "Invalid Auth data", err.Error())
			return
		case codes.PermissionDenied:
			common.Response(c, 403, "PERMISSION_DENIED", "Access denied", err.Error())
			return
		default:
			common.Response(c, 500, "INTERNAL_ERROR", "Internal server error", err.Error())
			return
		}
	}
	c.JSON(201, gin.H{
		"id":       resp.Id,
		"password": resp.Password,
	})
}
