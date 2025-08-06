package user

import (
	"context"

	"github.com/shinplay/ent"
	"github.com/shinplay/ent/user"
	"github.com/shinplay/internal/config"
	"go.uber.org/zap"
)

type UserRepositoryIntr interface {
	GetByPhoneNumber(ctx context.Context, phoneNumber string) (*ent.User, error)
	CreateByPhoneNumber(ctx context.Context, phoneNumber string) (*ent.User, error)
	FindByEmail(ctx context.Context, email string) (*ent.User, error)
	CreateByEmail(ctx context.Context, email string) (*ent.User, error)
	FindByUsername(ctx context.Context, username string) (*ent.User, error)
	UpdateUsername(ctx context.Context, userID string, newUsername string) error
}

type UserRepository struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{client: client}
}

// CreateByPhoneNumber implements UserRepository.
func (r *UserRepository) CreateByPhoneNumber(ctx context.Context, phoneNumber string) (*ent.User, error) {
	return r.client.User.Create().
		SetPhoneNumber(phoneNumber).
		Save(ctx)
}

// GetByPhoneNumber implements UserRepository.
func (r *UserRepository) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*ent.User, error) {
	config.GetConfig().Logger.Info("Getting user by phone number", zap.String("phoneNumber", phoneNumber))
	return r.client.User.Query().Where(user.PhoneNumberEQ(phoneNumber)).Only(ctx)
}

// FindByEmail implements UserRepository.
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*ent.User, error) {
	return r.client.User.Query().Where(user.EmailEQ(email)).Only(ctx)
}

// CreateByEmail implements UserRepository.
func (r *UserRepository) CreateByEmail(ctx context.Context, email string) (*ent.User, error) {
	return r.client.User.Create().
		SetEmail(email).
		Save(ctx)
}

// FindByUsername implements UserRepository.
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*ent.User, error) {
	return r.client.User.Query().Where(user.UsernameEQ(username)).Only(ctx)
}

// UpdateUsername implements UserRepository.
func (r *UserRepository) UpdateUsername(ctx context.Context, authId string, newUsername string) (*ent.User, error) {
	user, err := r.client.User.Query().Where(user.AuthIDEQ(authId)).Only(ctx)
	if err != nil {
		return nil, err
	}

	return r.client.User.UpdateOne(user).
		SetUsername(newUsername).
		Save(ctx)
}
