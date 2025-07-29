package db

import (
	"context"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/shinplay/ent"
	"github.com/shinplay/internal/config"
	"go.uber.org/zap"
)

func InitializeDatabase() *ent.Client {
	config := config.GetConfig()

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

	// defer client.Close()

	users, err := client.User.Query().All(context.Background())

	if err != nil {
		config.Logger.Fatal("failed to query users: %v", zap.Error(err))
	}

	if len(users) == 0 {
		config.Logger.Info("No users found in the database")
	}

	for _, user := range users {
		config.Logger.Info("User found %s", zap.Int("id", user.ID), zap.String("phone", user.PhoneNumber))
	}

	config.Logger.Info("PostgreSQL database connection established successfully :check_mark:")

	return client
}
