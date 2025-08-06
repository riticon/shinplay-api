package user

import (
	"context"

	"github.com/shinplay/ent"
	"github.com/shinplay/ent/user"
	"github.com/shinplay/internal/config"
	"go.uber.org/zap"
)

type UserServiceIntr interface {
	FindOrCreateByPhone(phoneNumber string) (*ent.User, error)
	FindOrCreateByEmail(email string) (*ent.User, error)
	FindByPhone(phoneNumber string) (*ent.User, error)
	FindByUsername(username string) (*ent.User, error)
	ChangeUsername(userID string, newUsername string) error
	FindUserByAuthID(authID string) (*ent.User, error)
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
			s.config.Logger.Info("Failed to get user by email", zap.Error(err))
			return nil, err
		}
	}

	return user, nil
}

func (s *UserService) FindByUsername(username string) (*ent.User, error) {
	user, err := s.userRepository.FindByUsername(s.ctx, username)
	if err != nil {
		s.config.Logger.Info("Failed to get user by username", zap.Error(err))
		return nil, err
	}

	return user, nil
}

func (s *UserService) ChangeUsername(userID string, newUsername string) (*ent.User, bool, error) {
	// check if the new username is already taken
	existingUser, err := s.userRepository.FindByUsername(s.ctx, newUsername)
	if err != nil && !ent.IsNotFound(err) {
		s.config.Logger.Error("Failed to check existing username", zap.Error(err))
		return nil, false, err
	}

	if existingUser != nil {
		s.config.Logger.Info("Username is already taken", zap.String("newUsername", newUsername))
		return nil, true, nil
	}

	user, err := s.userRepository.UpdateUsername(s.ctx, userID, newUsername)
	if err != nil {
		s.config.Logger.Error("Failed to change username", zap.String("userID", userID), zap.Error(err))
		return nil, false, err
	}

	s.config.Logger.Info("Username changed successfully", zap.String("userID", userID), zap.String("newUsername", newUsername))
	return user, false, nil
}

func (s *UserService) FindUserByAuthID(authID string) (*ent.User, error) {
	user, err := s.userRepository.client.User.Query().Where(user.AuthIDEQ(authID)).Only(s.ctx)
	if err != nil {
		s.config.Logger.Info("Failed to find user by auth ID", zap.String("authID", authID), zap.Error(err))
		return nil, err
	}

	return user, nil
}
