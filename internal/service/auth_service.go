package service

import (
	"time"

	"github.com/asikeida/xiaomajizhang/internal/auth"
	"github.com/asikeida/xiaomajizhang/internal/config"
	"github.com/asikeida/xiaomajizhang/internal/dto"
	apperrors "github.com/asikeida/xiaomajizhang/internal/errors"
	"github.com/asikeida/xiaomajizhang/internal/model"
	"github.com/asikeida/xiaomajizhang/internal/repository"
	"go.uber.org/zap"
)

type AuthService struct {
	cfg           config.JWTConfig
	users         *repository.UserRepository
	refreshTokens *repository.RefreshTokenRepository
	logger        *zap.Logger
}

func NewAuthService(cfg config.JWTConfig, users *repository.UserRepository, refreshTokens *repository.RefreshTokenRepository, logger *zap.Logger) *AuthService {
	return &AuthService{cfg: cfg, users: users, refreshTokens: refreshTokens, logger: logger}
}

func (s *AuthService) Register(req dto.RegisterRequest) (*dto.UserResponse, error) {
	exists, err := s.users.ExistsByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperrors.ErrUsernameExists
	}

	exists, err = s.users.ExistsByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperrors.ErrEmailExists
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	nickname := req.Nickname
	if nickname == "" {
		nickname = req.Username
	}

	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		Nickname:     nickname,
		Status:       1,
	}
	if err := s.users.Create(user); err != nil {
		return nil, err
	}

	s.logger.Info("user registered", zap.Uint64("user_id", user.ID), zap.String("username", user.Username))
	resp := toUserResponse(user)
	return &resp, nil
}

func (s *AuthService) Login(req dto.LoginRequest) (*dto.TokenResponse, error) {
	user, err := s.users.FindByAccount(req.Account)
	if err != nil {
		return nil, err
	}
	if user == nil || !auth.CheckPassword(req.Password, user.PasswordHash) {
		return nil, apperrors.ErrInvalidCredentials
	}
	if user.Status != 1 {
		return nil, apperrors.ErrUserDisabled
	}

	resp, err := s.issueTokens(user)
	if err != nil {
		return nil, err
	}
	if err := s.users.UpdateLastLoginAt(user.ID, time.Now()); err != nil {
		return nil, err
	}

	s.logger.Info("user login success", zap.Uint64("user_id", user.ID), zap.String("username", user.Username))
	return resp, nil
}

func (s *AuthService) Refresh(refreshToken string) (*dto.TokenResponse, error) {
	tokenHash := auth.HashToken(refreshToken)
	storedToken, err := s.refreshTokens.FindByHash(tokenHash)
	if err != nil {
		return nil, err
	}
	if storedToken == nil || storedToken.RevokedAt != nil || time.Now().After(storedToken.ExpiresAt) {
		return nil, apperrors.ErrRefreshTokenInvalid
	}

	user, err := s.users.FindByID(storedToken.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil || user.Status != 1 {
		return nil, apperrors.ErrRefreshTokenInvalid
	}

	return s.issueTokens(user)
}

func (s *AuthService) Logout(refreshToken string) error {
	tokenHash := auth.HashToken(refreshToken)
	storedToken, err := s.refreshTokens.FindByHash(tokenHash)
	if err != nil {
		return err
	}
	if storedToken == nil {
		return nil
	}
	return s.refreshTokens.RevokeByHash(tokenHash, time.Now())
}

func (s *AuthService) issueTokens(user *model.User) (*dto.TokenResponse, error) {
	accessToken, expiresIn, err := auth.GenerateAccessToken(s.cfg.Secret, user.ID, user.Username, s.cfg.AccessTokenExpiration)
	if err != nil {
		return nil, err
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}
	storedToken := &model.RefreshToken{
		UserID:    user.ID,
		TokenHash: auth.HashToken(refreshToken),
		ExpiresAt: time.Now().Add(s.cfg.RefreshTokenExpiration),
	}
	if err := s.refreshTokens.Create(storedToken); err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		User:         toUserResponse(user),
	}, nil
}

func toUserResponse(user *model.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		AvatarURL: user.AvatarURL,
	}
}
