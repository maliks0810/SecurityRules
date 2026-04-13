package snowflake

import (
	"database/sql"
	"fmt"

	sf "github.com/snowflakedb/gosnowflake"

	"securityrules/security-rules/configs"
	"securityrules/security-rules/internal/utils/log"
)

var DB *sql.DB

func Connect() {
	cfg := configs.EnvConfigs

	sfConfig := &sf.Config{
		Account:       cfg.SnowflakeAccount,
		User:          cfg.SnowflakeUser,
		Role:          cfg.SnowflakeRole,
		Database:      cfg.SnowflakeDatabase,
		Warehouse:     cfg.SnowflakeWarehouse,
		Schema:        cfg.SnowflakeSchema,
		Authenticator: sf.AuthTypeExternalBrowser,
	}

	dsn, err := sf.DSN(sfConfig)
	if err != nil {
		log.Logger.Fatal(fmt.Sprintf("Failed to build Snowflake DSN: %v", err))
	}

	DB, err = sql.Open("snowflake", dsn)
	if err != nil {
		log.Logger.Fatal(fmt.Sprintf("Failed to open Snowflake connection: %v", err))
	}

	if err := DB.Ping(); err != nil {
		log.Logger.Fatal(fmt.Sprintf("Failed to ping Snowflake: %v", err))
	}

	log.Logger.Info("Connected to Snowflake")
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
