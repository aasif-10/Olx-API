package main

import (
	"log"
	"os"

	"github.com/aasif-10/Olx-API/internals/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg := config.MustLoad()

	m, err := migrate.New(
		"file://migrations",
		cfg.DatabaseUrl)

	if err != nil {
		log.Fatalf("migrate.new: %v", err)
	}

	if len(os.Args) < 2 {
		log.Fatalf("usage: make migrate <up | down>")
	}

	switch os.Args[1] {
	case "up":
		err := m.Up()
		if err != nil {
			log.Fatalf("m.up: %v", err)
		}

	case "down":
		err := m.Steps(-1)
		if err != nil {
			log.Fatalf("m.down: %v", err)
		}

	default:
		log.Fatalf("unknown command %s", os.Args[1])
	}
}
