package config

import (
	"log"
	"os"
	"time"
)

type Configuration struct {
	DatabaseName        string
	DatabaseHost        string
	DatabaseUser        string
	DatabasePassword    string
	DatabasePath        string
	MigrateToVersion    string
	MigrationLocation   string
	FileStorageLocation string
	JwtSecret           string
	JwtTTL              time.Duration
	MonobankPrivateKey  string // Токен з особистого кабінету https://web.monobank.ua/ або тестовий токен з https://api.monobank.ua/
}

func GetConfiguration() Configuration {
	return Configuration{
		DatabaseName:        getOrFail("DB_NAME"),
		DatabaseHost:        getOrFail("DB_HOST"),
		DatabaseUser:        getOrFail("DB_USER"),
		DatabasePassword:    getOrFail("DB_PASSWORD"),
		MigrateToVersion:    getOrDefault("MIGRATE", "latest"),
		MigrationLocation:   getOrDefault("MIGRATION_LOCATION", "internal/infra/database/migrations"),
		FileStorageLocation: getOrDefault("FILES_LOCATION", "file_storage"),
		JwtSecret:           getOrDefault("JWT_SECRET", "1234567890"),
		JwtTTL:              72 * time.Hour,
		MonobankPrivateKey:  getOrDefault("MONOBANK_PRIVATE_KEY", "uES2_x-N_rd3eysY_-SXsoqBAIgmK4lLnqZpRZMKdAM4"),
	}
}

//nolint:unused
func getOrFail(key string) string {
	env, set := os.LookupEnv(key)
	if !set || env == "" {
		log.Fatalf("%s env var is missing", key)
	}
	return env
}

func getOrDefault(key, defaultVal string) string {
	env, set := os.LookupEnv(key)
	if !set {
		return defaultVal
	}
	return env
}
