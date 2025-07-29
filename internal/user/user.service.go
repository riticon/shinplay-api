package user

import (
	"context"

	"github.com/shinplay/ent"
	"github.com/shinplay/internal/config"
	"go.uber.org/zap"
)

type UserService interface {
	FindOrCreateByPhone(phoneNumber string) (*ent.User, error)
}

// UserService provides methods to manage user-related operations.
type UserServiceImpl struct {
	userRepository UserRepository
	config         *config.Config
}

// NewUserService creates a new UserService instance.
func NewUserService(userRepository UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		userRepository: userRepository,
		config:         config.GetConfig(),
	}
}

func (s *UserServiceImpl) FindOrCreateByPhone(phoneNumber string) (*ent.User, error) {
	s.config.Logger.Info("Finding or creating user by phone number", zap.String("phoneNumber", phoneNumber))
	user, err := s.userRepository.GetByPhoneNumber(context.Background(), phoneNumber)

	if err != nil {
		if ent.IsNotFound(err) {
			user, err = s.userRepository.CreateByPhoneNumber(context.Background(), phoneNumber)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return user, nil
}
