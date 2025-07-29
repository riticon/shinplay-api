package auth

import (
	"context"

	"github.com/shinplay/ent"
	"github.com/shinplay/internal/config"
	"go.uber.org/zap"
)

type OTPService interface {
	CreateNewOTP(user *ent.User) (*ent.OTP, error)
}

type OTPServiceImpl struct {
	otpRepository OTPRepository
	config        *config.Config
}

func NewOTPService(otpRepository OTPRepository) *OTPServiceImpl {
	return &OTPServiceImpl{
		otpRepository: otpRepository,
		config:        config.GetConfig(),
	}
}

func (s *OTPServiceImpl) CreateNewOTP(user *ent.User) (*ent.OTP, error) {
	otp, err := s.otpRepository.CreateNewOTP(context.Background(), user)

	if err != nil {
		s.config.Logger.Error("Failed to create new OTP", zap.Error(err))
		return nil, err
	}

	return otp, nil
}
