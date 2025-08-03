package config

import (
	"sync"

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

type WhatsAppConfig struct {
	Token   string
	PhoneId string
}

type Config struct {
	Name        string
	Environment string
	Database    DatabaseConfig
	Server      ServerConfig
	WhatsApp    WhatsAppConfig
	JWTSecret   string
	Logger      *zap.Logger
}

// singleton instance of Config.
var (
	instance *Config
	once     sync.Once
)

// GetConfig returns the singleton instance of Config.
func GetConfig() *Config {

	once.Do(func() {
		env := LoadEnv()

		instance = &Config{
			Name:        "default",
			Environment: env.Environment,
			JWTSecret:   env.JWTSecret,
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
			WhatsApp: WhatsAppConfig{
				Token:   env.WhatsAppToken,
				PhoneId: env.WhatsAppPhoneId,
			},
			Logger: nil,
		}
		instance.InitalizeLogger()

		instance.Logger.Info("Config initialized", zap.String("environment", instance.Environment))
	})

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
