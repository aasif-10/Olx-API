package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	Env  string
}

func MustLoad() Config {
	godotenv.Load()

	port := os.Getenv("PORT")
	env := os.Getenv("ENV")

	if port == "" {
		panic("PORT is required")
	}

	if env == "" {
		panic("ENV is required")
	}

	return Config{Port: port, Env: env}
}
