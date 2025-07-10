package repository

import (
	"context"

	"github.com/shinplay/ent"
	"github.com/shinplay/ent/user"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int) (*ent.User, error)
	GetByEmail(ctx context.Context, email string) (*ent.User, error)
	Create(ctx context.Context, name, email string) (*ent.User, error)
	UpdateName(ctx context.Context, id int, name string) error
	Delete(ctx context.Context, id int) error
}

type userRepository struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) UserRepository {
	return &userRepository{client: client}
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*ent.User, error) {
	return r.client.User.Get(ctx, id)
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*ent.User, error) {
	return r.client.User.Query().Where(user.Email(email)).Only(ctx)
}

func (r *userRepository) Create(ctx context.Context, firstName, email string) (*ent.User, error) {
	return r.client.User.Create().SetFirstName(firstName).SetEmail(email).Save(ctx)
}

func (r *userRepository) UpdateName(ctx context.Context, id int, lastName string) error {
	return r.client.User.UpdateOneID(id).SetLastName(lastName).Exec(ctx)
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	return r.client.User.DeleteOneID(id).Exec(ctx)
}
