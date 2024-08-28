package main

import (
	"auth/internal/config"
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

func main() {
	var configPath, migrationsPath string
	flag.StringVar(&configPath, "config", "", "The path to the configuration file")
	flag.StringVar(&migrationsPath, "migrations", "", "The path to migration")
	flag.Parse()

	if configPath == "" {
		log.Fatal("You need to specify the path to the configuration")
	}
	if migrationsPath == "" {
		log.Fatal("You need to specify the path to migrations")
	}

	cfg := config.MustLoad(configPath)

	var dsn string

	switch cfg.Storage.Type {
	case "sqlite":
		dsn = fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", cfg.Storage.Sqlite.StoragePath, "migrations")
	case "postgres":
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s&x-migrations-table=%s",
			cfg.Storage.Postgres.Username,
			cfg.Storage.Postgres.Password,
			cfg.Storage.Postgres.Host,
			cfg.Storage.Postgres.Port,
			cfg.Storage.Postgres.Database,
			cfg.Storage.Postgres.Sslmode,
			"migrations")
	default:
		log.Fatalf("Unsupported storage type: %s", cfg.Storage.Type)
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("There are no new migrations")
			return
		}
		log.Fatal(err)
	}
	fmt.Println("Migrations have been successfully applied")
}
