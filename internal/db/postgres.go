package db

import (
	"fmt"

	_ "github.com/lib/pq"
	"github.com/shinplay/ent"
	"github.com/shinplay/internal/config"
)

func NewPostgresDB() (*ent.Client, error) {
	config := config.GetConfig()

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
		return nil, err
	}

	return db, nil
}
