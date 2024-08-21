package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	HTTPPort string

	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDatabase string

	KafkaBrokers []string

	DefaultOffset string
	DefaultLimit  string

	TokenKey string
}

// Load ...
func Load() Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found:", err)
	}

	config := Config{}

	config.HTTPPort = cast.ToString(getOrReturnDefaultValue("HTTP_PORT", ":8085"))

	config.PostgresHost = cast.ToString(getOrReturnDefaultValue("POSTGRES_HOST", "postgres-pr"))
	config.PostgresPort = cast.ToInt(getOrReturnDefaultValue("POSTGRES_PORT", 5835))
	config.PostgresUser = cast.ToString(getOrReturnDefaultValue("POSTGRES_USER", "postgres"))
	config.PostgresPassword = cast.ToString(getOrReturnDefaultValue("POSTGRES_PASSWORD", "1234"))
	config.PostgresDatabase = cast.ToString(getOrReturnDefaultValue("POSTGRES_DATABASE", "product"))

	config.KafkaBrokers = parseKafkaBrokers(getOrReturnDefaultValue("KAFKA_BROKERS", "kafka:9092"))

	config.DefaultOffset = cast.ToString(getOrReturnDefaultValue("DEFAULT_OFFSET", "0"))
	config.DefaultLimit = cast.ToString(getOrReturnDefaultValue("DEFAULT_LIMIT", "10"))

	config.TokenKey = cast.ToString(getOrReturnDefaultValue("TokenKey", "my_secret_key"))

	return config
}

func getOrReturnDefaultValue(key string, defaultValue interface{}) interface{} {
	val, exists := os.LookupEnv(key)

	if exists {
		return val
	}

	return defaultValue
}

func parseKafkaBrokers(brokers interface{}) []string {
	switch v := brokers.(type) {
	case string:
		return strings.Split(v, ",")
	case []string:
		return v
	default:
		return []string{"kafka:9092"}
	}
}
