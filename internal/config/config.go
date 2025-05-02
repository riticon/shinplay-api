package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Name   string
	Env    Env
	Logger *zap.Logger
}

// singleton instance of Config
var instance *Config

// GetConfig returns the singleton instance of Config
func GetConfig() *Config {
	if instance == nil {
		instance = &Config{
			Name: "default",
			Env:  LoadEnv(),
		}

		instance.InitalizeLogger()
	}

	return instance
}

func (c *Config) InitalizeLogger() {
	logger, _ := zap.NewProduction()

	if c.IsDevelopment() {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, _ = config.Build()
	}

	defer logger.Sync()
	c.Logger = logger
}

func (c *Config) IsDevelopment() bool {
	return c.Env.Environment == "development"
}

func (c *Config) IsProduction() bool {
	return c.Env.Environment == "production"
}

func (c *Config) IsStaging() bool {
	return c.Env.Environment == "staging"
}
