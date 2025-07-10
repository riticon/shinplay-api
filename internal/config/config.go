package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	Name     string
}

type ServerConfig struct {
	Port string
	Host string
	CORS string
}

type Config struct {
	Name        string
	Environment string
	Database    DatabaseConfig
	Server      ServerConfig
	Logger      *zap.Logger
}

// singleton instance of Config.
var instance *Config

// GetConfig returns the singleton instance of Config.
func GetConfig() *Config {
	if instance == nil {
		print("Initializing Config...")

		env := LoadEnv()
		instance = &Config{
			Name:        "default",
			Environment: env.Environment,
			Database: DatabaseConfig{
				Host:     env.DBHost,
				User:     env.DBUser,
				Password: env.DBPassword,
				Name:     env.DBName,
			},
			Server: ServerConfig{
				Port: env.ServerPort,
				Host: env.ServerHost,
				CORS: env.CORS,
			},
			Logger: nil,
		}

		instance.InitalizeLogger()
	}

	return instance
}

func (c *Config) InitalizeLogger() {
	logger, _ := zap.NewProduction()

	if c.IsDevelopment() {
		logConfig := zap.NewDevelopmentConfig()
		logConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, _ = logConfig.Build()
	}

	defer logger.Sync() //nolint:errcheck
	c.Logger = logger
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

func (c *Config) IsStaging() bool {
	return c.Environment == "staging"
}
