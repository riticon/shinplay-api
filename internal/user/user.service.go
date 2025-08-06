package user

import (
	"context"

	"github.com/shinplay/ent"
	"github.com/shinplay/internal/config"
	"go.uber.org/zap"
)

type UserServiceIntr interface {
	FindOrCreateByPhone(phoneNumber string) (*ent.User, error)
	FindOrCreateByEmail(email string) (*ent.User, error)
	FindByPhone(phoneNumber string) (*ent.User, error)
}

// UserService provides methods to manage user-related operations.
type UserService struct {
	ctx            context.Context
	userRepository *UserRepository
	config         *config.Config
}

// NewUserService creates a new UserService instance.
func NewUserService(userRepository *UserRepository, config *config.Config, ctx context.Context) *UserService {
	return &UserService{
		config:         config,
		ctx:            ctx,
		userRepository: userRepository,
	}
}

func (s *UserService) FindOrCreateByPhone(phoneNumber string) (*ent.User, error) {
	user, err := s.userRepository.GetByPhoneNumber(context.Background(), phoneNumber)

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

func (s *UserService) FindByPhone(phoneNumber string) (*ent.User, error) {
	user, err := s.userRepository.GetByPhoneNumber(context.Background(), phoneNumber)
	if err != nil {
		s.config.Logger.Error("Failed to get user by phone number", zap.Error(err))
		return nil, err
	}

	return user, nil
}

func (s *UserService) FindOrCreateByEmail(email string) (*ent.User, error) {
	user, err := s.userRepository.FindByEmail(context.Background(), email)

	if err != nil {
		if ent.IsNotFound(err) {
			s.config.Logger.Info("User not found, creating new user", zap.String("email", email))
			user, err = s.userRepository.CreateByEmail(context.Background(), email)
			if err != nil {
				s.config.Logger.Error("Failed to create user by email", zap.Error(err))
				return nil, err
			}
		} else {
			s.config.Logger.Error("Failed to get user by email", zap.Error(err))
			return nil, err
		}
	}

	return user, nil
}
