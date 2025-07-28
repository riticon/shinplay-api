package user

import (
	"context"

	"github.com/shinplay/ent"
	"github.com/shinplay/internal/config"
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
func NewUserService() *UserServiceImpl {
	return &UserServiceImpl{
		userRepository: NewUserRepository(),
		config:         config.GetConfig(),
	}
}

func (s *UserServiceImpl) FindOrCreateByPhone(phoneNumber string) (*ent.User, error) {
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
