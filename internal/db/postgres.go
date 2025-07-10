package db

import (
	"database/sql"
	"fmt"

	"github.com/shinplay/internal/config"
)

func NewPostgresDB() (*sql.DB, error) {
	config := config.GetConfig()

	db, err := sql.Open(
		"postgres",
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
			config.Database.Host,
			config.Database.User,
			config.Database.Name),
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}
