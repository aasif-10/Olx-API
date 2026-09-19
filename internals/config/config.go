package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                string
	Env                 string
	DatabaseUrl         string
	JwtKey              string
	StorageProjectId    string
	StorageAccessKey    string
	StorageAccessSecret string
	StorageBucket       string
}

func MustLoad() Config {
	godotenv.Load()

	port := os.Getenv("PORT")
	env := os.Getenv("ENV")
	dbUrl := os.Getenv("DATABASE_URL")
	jwtKey := os.Getenv("JWT_KEY")
	storageProjectId := os.Getenv("STORAGE_PROJECT_ID")
	storageAccessKey := os.Getenv("STORAGE_ACCESS_KEY")
	storageAccountSecret := os.Getenv("STORAGE_ACCESS_SECRET")
	storageBucket := os.Getenv("STORAGE_BUCKET")

	if port == "" {
		panic("PORT is required")
	}

	if env == "" {
		panic("ENV is required")
	}

	if dbUrl == "" {
		panic("DATABASE URL is required")
	}

	if jwtKey == "" {
		panic("JWT_KEY is required")
	}

	if storageProjectId == "" {
		panic("STORAGE_PROJECT_ID is required")
	}

	if storageAccessKey == "" {
		panic("STORAGE_ACCESS_KEY is required")
	}

	if storageAccountSecret == "" {
		panic("STORAGE_ACCESS_SECRET is required")
	}

	if storageBucket == "" {
		panic("STORAGE_BUCKET is required")
	}

	return Config{
		Port:                port,
		Env:                 env,
		DatabaseUrl:         dbUrl,
		JwtKey:              jwtKey,
		StorageProjectId:    storageProjectId,
		StorageAccessKey:    storageAccessKey,
		StorageAccessSecret: storageAccountSecret,
		StorageBucket:       storageBucket}
}
