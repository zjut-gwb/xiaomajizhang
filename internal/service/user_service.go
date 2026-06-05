package service

import (
	"github.com/asikeida/xiaomajizhang/internal/dto"
	apperrors "github.com/asikeida/xiaomajizhang/internal/errors"
	"github.com/asikeida/xiaomajizhang/internal/repository"
)

type UserService struct {
	users *repository.UserRepository
}

func NewUserService(users *repository.UserRepository) *UserService {
	return &UserService{users: users}
}

func (s *UserService) GetByID(userID uint64) (*dto.UserResponse, error) {
	user, err := s.users.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperrors.ErrUserNotFound
	}
	resp := dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		AvatarURL: user.AvatarURL,
	}
	return &resp, nil
}
