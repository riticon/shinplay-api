package user

import (
	"context"

	"github.com/shinplay/ent"
	"github.com/shinplay/internal/config"
	"go.uber.org/zap"
)

type UserServiceIntr interface {
	FindOrCreateByPhone(phoneNumber string) (*ent.User, error)
}

// UserService provides methods to manage user-related operations.
type UserService struct {
	userRepository *UserRepository
	config         *config.Config
}

// NewUserService creates a new UserService instance.
func NewUserService(userRepository *UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
		config:         config.GetConfig(),
	}
}

func (s *UserService) FindOrCreateByPhone(phoneNumber string) (*ent.User, error) {
	s.config.Logger.Info("Finding or creating user by phone number", zap.String("phoneNumber", phoneNumber))
	user, err := s.userRepository.GetByPhoneNumber(context.Background(), phoneNumber)
	s.config.Logger.Info("User found", zap.Any("user", user))

	if err != nil {
		if ent.IsNotFound(err) {
			s.config.Logger.Info("User not found, creating new user", zap.String("phoneNumber", phoneNumber))
			user, err = s.userRepository.CreateByPhoneNumber(context.Background(), phoneNumber)
			if err != nil {
				s.config.Logger.Error("Failed to create user by phone number", zap.Error(err))
				return nil, err
			}
		} else {
			s.config.Logger.Error("Failed to get user by phone number", zap.Error(err))
			return nil, err
		}
	}

	return user, nil
}
