package config

import "os"

type Config struct {
	Port   string
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
}

type SnowflakeConfig struct {
	Account       string
	User          string
	Role          string
	Database      string
	Warehouse     string
	Schema        string
	Authenticator string
}

func Load() *Config {
	return &Config{
		Port:   getEnv("PORT", "3001"),
		DBHost: getEnv("DB_HOST", "localhost"),
		DBPort: getEnv("DB_PORT", "5432"),
		DBUser: getEnv("DB_USER", "postgres"),
		DBPass: getEnv("DB_PASSWORD", "1010data"),
		DBName: getEnv("DB_NAME", "DATA_QUALITY"),
	}
}

func LoadSnowflake() *SnowflakeConfig {
	return &SnowflakeConfig{
		Account:       getEnv("SNOWFLAKE_ACCOUNT", ""),
		User:          getEnv("SNOWFLAKE_USER", ""),
		Role:          getEnv("SNOWFLAKE_ROLE", ""),
		Database:      getEnv("SNOWFLAKE_DATABASE", ""),
		Warehouse:     getEnv("SNOWFLAKE_WAREHOUSE", ""),
		Schema:        getEnv("SNOWFLAKE_SCHEMA", "PUBLIC"),
		Authenticator: getEnv("SNOWFLAKE_AUTHENTICATOR", "externalbrowser"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
