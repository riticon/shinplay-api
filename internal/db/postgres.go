package db

import (
	"context"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/shinplay/ent"
	"github.com/shinplay/internal/config"
	"go.uber.org/zap"
)

func InitializeDatabase(config *config.Config) *ent.Client {
	config.Logger.Info("Initializing PostgreSQL database connection")

	client, err := ent.Open(
		"postgres",
		fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			config.Database.Host,
			5432, // Default PostgreSQL port
			config.Database.User,
			config.Database.Password,
			config.Database.Name),
	)

	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		config.Logger.Fatal("failed creating schema resources: %v", zap.Error(err))
	}

	if err != nil {
		config.Logger.Fatal("failed to query users: %v", zap.Error(err))
	}

	config.Logger.Info("PostgreSQL database connection established successfully")

	return client
}
