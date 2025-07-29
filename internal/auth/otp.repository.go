package auth

import (
	"context"

	"github.com/shinplay/ent"
)

type OTPRepositoryIntr interface {
	CreateNewOTP(ctx context.Context, user *ent.User) (*ent.OTP, error)
}

type OTPRepository struct {
	client *ent.Client
}

func NewOTPRepository(client *ent.Client) *OTPRepository {
	return &OTPRepository{client: client}
}

// CreateNewOTP implements OTPRepository.
func (o *OTPRepository) CreateNewOTP(ctx context.Context, user *ent.User) (*ent.OTP, error) {
	return o.client.OTP.Create().
		SetUser(user).
		Save(ctx)
}
