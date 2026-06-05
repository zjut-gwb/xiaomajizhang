package handler

import (
	"errors"

	"github.com/asikeida/xiaomajizhang/internal/dto"
	apperrors "github.com/asikeida/xiaomajizhang/internal/errors"
	"github.com/asikeida/xiaomajizhang/internal/response"
	"github.com/asikeida/xiaomajizhang/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	auth   *service.AuthService
	logger *zap.Logger
}

func NewAuthHandler(auth *service.AuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{auth: auth, logger: logger}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, 40000, "请求参数错误")
		return
	}

	user, err := h.auth.Register(req)
	if err != nil {
		h.writeAuthError(c, err)
		return
	}
	response.OK(c, user)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, 40000, "请求参数错误")
		return
	}

	tokens, err := h.auth.Login(req)
	if err != nil {
		h.writeAuthError(c, err)
		return
	}
	response.OK(c, tokens)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, 40000, "请求参数错误")
		return
	}

	tokens, err := h.auth.Refresh(req.RefreshToken)
	if err != nil {
		h.writeAuthError(c, err)
		return
	}
	response.OK(c, tokens)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, 40000, "请求参数错误")
		return
	}
	if err := h.auth.Logout(req.RefreshToken); err != nil {
		h.logger.Error("logout failed", zap.Error(err))
		response.Error(c, 500, 50000, "服务器内部错误")
		return
	}
	response.OK(c, gin.H{"success": true})
}

func (h *AuthHandler) writeAuthError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, apperrors.ErrUsernameExists):
		response.Error(c, 409, 40900, "用户名已存在")
	case errors.Is(err, apperrors.ErrEmailExists):
		response.Error(c, 409, 40900, "邮箱已存在")
	case errors.Is(err, apperrors.ErrInvalidCredentials):
		response.Error(c, 401, 40100, "用户名或密码错误")
	case errors.Is(err, apperrors.ErrUserDisabled):
		response.Error(c, 403, 40300, "用户已禁用")
	case errors.Is(err, apperrors.ErrRefreshTokenInvalid):
		response.Error(c, 401, 40100, "刷新令牌无效")
	default:
		h.logger.Error("auth request failed", zap.Error(err))
		response.Error(c, 500, 50000, "服务器内部错误")
	}
}
