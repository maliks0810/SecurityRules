package database

import (
	"SecurityRules/config"
	"database/sql"
	"fmt"
	"log"

	sf "github.com/snowflakedb/gosnowflake"
)

var SnowflakeDB *sql.DB

func ConnectSnowflake(cfg *config.SnowflakeConfig) {
	sfConfig := &sf.Config{
		Account:       cfg.Account,
		User:          cfg.User,
		Role:          cfg.Role,
		Database:      cfg.Database,
		Warehouse:     cfg.Warehouse,
		Schema:        cfg.Schema,
		Authenticator: sf.AuthTypeExternalBrowser,
	}

	dsn, err := sf.DSN(sfConfig)
	if err != nil {
		log.Fatalf("Failed to build Snowflake DSN: %v", err)
	}

	SnowflakeDB, err = sql.Open("snowflake", dsn)
	if err != nil {
		log.Fatalf("Failed to open Snowflake connection: %v", err)
	}

	if err := SnowflakeDB.Ping(); err != nil {
		log.Fatalf("Failed to ping Snowflake: %v", err)
	}

	fmt.Println("Connected to Snowflake")
}

func CloseSnowflake() {
	if SnowflakeDB != nil {
		SnowflakeDB.Close()
	}
}
