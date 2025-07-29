package auth

import (
	"context"

	"github.com/shinplay/ent"
)

type OTPRepository interface {
	CreateNewOTP(ctx context.Context, user *ent.User) (*ent.OTP, error)
}

type otpRepository struct {
	client *ent.Client
}

func NewOTPRepository(client *ent.Client) OTPRepository {
	return &otpRepository{client: client}
}

// CreateNewOTP implements OTPRepository.
func (o *otpRepository) CreateNewOTP(ctx context.Context, user *ent.User) (*ent.OTP, error) {
	return o.client.OTP.Create().
		SetUser(user).
		Save(ctx)
}
