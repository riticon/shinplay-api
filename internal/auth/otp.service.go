package auth

import (
	"context"

	"github.com/shinplay/ent"
	"github.com/shinplay/internal/config"
	"go.uber.org/zap"
)

type OTPServiceIntr interface {
	CreateNewOTP(user *ent.User) (*ent.OTP, error)
}

type OTPService struct {
	otpRepository *OTPRepository
	config        *config.Config
}

func NewOTPService(otpRepository *OTPRepository) *OTPService {
	return &OTPService{
		otpRepository: otpRepository,
		config:        config.GetConfig(),
	}
}

func (s *OTPService) CreateNewOTP(user *ent.User) (*ent.OTP, error) {
	otp, err := s.otpRepository.CreateNewOTP(context.Background(), user)

	if err != nil {
		s.config.Logger.Error("Failed to create new OTP", zap.Error(err))
		return nil, err
	}

	return otp, nil
}
