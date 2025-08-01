package main

import (
	"errors"
	"fmt"
	"url-shortener/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("local.env"); err != nil {
		fmt.Errorf("could not find .env file: %s", err)
	}
	cfg := config.MustLoad()
	m, err := migrate.New(
		"file://"+cfg.MigrationsPath,
		cfg.StoragePath,
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		panic(err)
	}
	fmt.Println("migrations applied successfully")
}
