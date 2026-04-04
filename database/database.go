package database

import (
	"SecurityRules/config"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect(cfg *config.Config) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)

	var err error
	Pool, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	if err := Pool.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}

	log.Println("Connected to PostgreSQL")
}

func Migrate() {
	query := `
	CREATE TABLE IF NOT EXISTS security_rules (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT DEFAULT '',
		severity VARCHAR(20) NOT NULL CHECK (severity IN ('low', 'medium', 'high', 'critical')),
		category VARCHAR(255) DEFAULT '',
		enabled BOOLEAN DEFAULT false,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);`

	_, err := Pool.Exec(context.Background(), query)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Database migration complete")
}

func Close() {
	if Pool != nil {
		Pool.Close()
	}
}
