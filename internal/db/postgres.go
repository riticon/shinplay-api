package db

import (
	"context"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/shinplay/ent"
	"github.com/shinplay/internal/config"
	"go.uber.org/zap"
)

func InitializeDatabase() {
	config := config.GetConfig()

	config.Logger.Info("Initializing PostgreSQL database connection")

	db, err := ent.Open(
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

	if err := db.Schema.Create(context.Background()); err != nil {
		config.Logger.Fatal("failed creating schema resources: %v", zap.Error(err))
	}

	config.Logger.Info("PostgreSQL database connection established successfully :check_mark:")
}
