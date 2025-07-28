package user

import (
	"context"

	"github.com/shinplay/ent"
	"github.com/shinplay/ent/user"
)

type UserRepository interface {
	GetByPhoneNumber(ctx context.Context, phoneNumber string) (*ent.User, error)
	CreateByPhoneNumber(ctx context.Context, phoneNumber string) (*ent.User, error)
}

type userRepository struct {
	client *ent.Client
}

func NewUserRepository() UserRepository {
	return &userRepository{client: ent.NewClient()}
}

// CreateByPhoneNumber implements UserRepository.
func (r *userRepository) CreateByPhoneNumber(ctx context.Context, phoneNumber string) (*ent.User, error) {
	return r.client.User.Create().
		SetPhoneNumber(phoneNumber).
		SetAuthID(ent.User{}.AuthID).
		Save(ctx)
}

// GetByPhoneNumber implements UserRepository.
func (r *userRepository) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*ent.User, error) {
	return r.client.User.Query().Where(user.PhoneNumberEQ(phoneNumber)).Only(ctx)
}
