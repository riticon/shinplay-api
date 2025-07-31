package auth

import (
	"context"

	"github.com/shinplay/ent"
	"github.com/shinplay/ent/otp"
	"github.com/shinplay/ent/user"
)

type OTPRepositoryIntr interface {
	CreateNewOTP(ctx context.Context, user *ent.User) (*ent.OTP, error)
	FindOTPByUser(ctx context.Context, user *ent.User) (*ent.OTP, error)
	DeleteOTP(ctx context.Context, otpCode string, user *ent.User) error
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

func (o *OTPRepository) FindOTPByUser(ctx context.Context, otp_code string, user_ *ent.User) (*ent.OTP, error) {
	return o.client.OTP.Query().
		Where(otp.OtpEQ(otp_code)).
		Where(otp.HasUserWith(user.IDEQ(user_.ID))).
		First(ctx)
}

func (o *OTPRepository) DeleteOTP(ctx context.Context, otpCode string, user_ *ent.User) (int, error) {
	otpId, err := o.client.OTP.Delete().
		Where(otp.OtpEQ(otpCode)).
		Where(otp.HasUserWith(user.IDEQ(user_.ID))).
		Exec(ctx)

	if err != nil {
		return 0, err
	}

	return otpId, nil
}
