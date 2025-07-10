package service

import (
	"context"

	"github.com/shinplay/ent"
	"github.com/shinplay/internal/repository"
)

type UserService interface {
	GetByID(ctx context.Context, id int) (*ent.User, error)
	GetByEmail(ctx context.Context, email string) (*ent.User, error)
	Create(ctx context.Context, name, email string) (*ent.User, error)
	UpdateName(ctx context.Context, id int, name string) error
	Delete(ctx context.Context, id int) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetByID(ctx context.Context, id int) (*ent.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*ent.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

func (s *userService) Create(ctx context.Context, name, email string) (*ent.User, error) {
	return s.userRepo.Create(ctx, name, email)
}

func (s *userService) UpdateName(ctx context.Context, id int, name string) error {
	return s.userRepo.UpdateName(ctx, id, name)
}

func (s *userService) Delete(ctx context.Context, id int) error {
	return s.userRepo.Delete(ctx, id)
}
