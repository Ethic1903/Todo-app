package main

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"url-shortener/internal/config"
)

func main() {
	//var storagePath, migrationsPath, migrationsTable string
	//
	//flag.StringVar(&storagePath, "storage-path", "", "path to storage")
	//flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	//flag.StringVar(&migrationsTable, "migrations-table", "", "name of migrations")
	//flag.Parse()
	//
	//if storagePath == "" {
	//	panic("storage-path is required")
	//}
	//if migrationsPath == "" {
	//	panic("migrations-path is required")
	//}
	if err := godotenv.Load("local.env"); err != nil {
		fmt.Errorf("could not find .env file: %s", err)
	}
	cfg := config.MustLoad()
	m, err := migrate.New(
		"file://"+cfg.MigrationsPath,
		/* fmt.Sprintf("postgres://%s?x-migrations-table=%s",*/ cfg.StoragePath)
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
