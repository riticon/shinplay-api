package otp

import (
	"context"

	"github.com/shinplay/ent"
	"github.com/shinplay/internal/config"
	"go.uber.org/zap"
)

type OTPServiceIntr interface {
	CreateNewOTP(user *ent.User) (*ent.OTP, error)
	IsOTPValid(otp string, user *ent.User) (bool, error)
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

func (s *OTPService) IsOTPValid(otpCode string, user *ent.User) (bool, error) {
	s.config.Logger.Info("Validating OTP", zap.String("otpCode", otpCode), zap.Int("userId", user.ID))

	otp, err := s.otpRepository.FindOTPByUser(context.Background(), otpCode, user)

	if err != nil {
		s.config.Logger.Debug("Failed to find OTP for user", zap.Error(err))
		return false, err
	}

	if otp == nil {
		s.config.Logger.Info("No OTP found for user", zap.Int("userId", user.ID))
		return false, nil
	}

	s.config.Logger.Info("OTP found for user", zap.Int("userId", user.ID), zap.String("otpCode", otp.Otp))
	return true, nil
}

func (s *OTPService) ExpireOtp(otpCode string, userId int) (bool, error) {
	s.config.Logger.Info("Expiring OTP", zap.String("otpCode", otpCode), zap.Int("userId", userId))

	deletedId, err := s.otpRepository.DeleteOTP(context.Background(), otpCode, userId)

	s.config.Logger.Info("OTP is used - Deleting", zap.Int("deletedId", deletedId))

	if err != nil {
		s.config.Logger.Error("Failed to delete OTP", zap.Error(err))
		return false, err
	}

	s.config.Logger.Info("OTP deleted successfully", zap.Int("deletedId", deletedId))
	return true, nil
}
