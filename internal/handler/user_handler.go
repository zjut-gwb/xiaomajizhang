package handler

import (
	"errors"

	apperrors "github.com/asikeida/xiaomajizhang/internal/errors"
	"github.com/asikeida/xiaomajizhang/internal/response"
	"github.com/asikeida/xiaomajizhang/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	users  *service.UserService
	logger *zap.Logger
}

func NewUserHandler(users *service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{users: users, logger: logger}
}

func (h *UserHandler) Me(c *gin.Context) {
	userIDValue, ok := c.Get("user_id")
	if !ok {
		response.Error(c, 401, 40100, "未登录")
		return
	}

	userID, ok := userIDValue.(uint64)
	if !ok {
		response.Error(c, 401, 40100, "未登录")
		return
	}

	user, err := h.users.GetByID(userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			response.Error(c, 404, 40400, "用户不存在")
			return
		}
		h.logger.Error("get current user failed", zap.Error(err), zap.Uint64("user_id", userID))
		response.Error(c, 500, 50000, "服务器内部错误")
		return
	}
	response.OK(c, user)
}
