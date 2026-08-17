package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Env         string
	DatabaseUrl string
}

func MustLoad() Config {
	godotenv.Load()

	port := os.Getenv("PORT")
	env := os.Getenv("ENV")
	dbUrl := os.Getenv("DATABASE_URL")

	if port == "" {
		panic("PORT is required")
	}

	if env == "" {
		panic("ENV is required")
	}

	if dbUrl == "" {
		panic("DATABASE URL is required")
	}

	return Config{Port: port, Env: env, DatabaseUrl: dbUrl}
}
