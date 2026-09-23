package config

import (
	"log"
	"net/url"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresDSN string
	RedisAddr   string
	RedisDB     int
	HTTPPort    string
}

const defaultDSN = "postgres://postgres:123@postgres:5432/em-db?sslmode=disable"

func Load() *Config {
	// .env необязателен — при docker-compose переменные придут из env
	_ = godotenv.Load()

	def, err := url.Parse(defaultDSN)
	if err != nil {
		log.Fatalf("cannot parse default dsn: %v", err)
	}
	defUser := def.User.Username()
	defPass, _ := def.User.Password()
	defHost := def.Hostname()
	defPort := def.Port()
	defDB := def.Path[1:] 

	host := getEnv("POSTGRES_HOST", defHost)
	port := getEnv("POSTGRES_PORT", defPort)
	user := getEnv("POSTGRES_USER", defUser)
	pass := getEnv("POSTGRES_PASSWORD", defPass)
	name := getEnv("POSTGRES_DB", defDB)

	dsn := "postgres://" + user + ":" + pass + "@" + host + ":" + port + "/" + name + "?sslmode=disable"

	return &Config{
		PostgresDSN: dsn,
		RedisAddr:   getEnv("REDIS_HOST", "redis") + ":" + getEnv("REDIS_PORT", "6379"),
		RedisDB:     getEnvInt("REDIS_DB", 0),
		HTTPPort:    getEnv("HTTP_PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}