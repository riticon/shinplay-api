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
