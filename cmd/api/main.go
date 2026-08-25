package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aasif-10/Olx-API/internals/config"
	"github.com/aasif-10/Olx-API/internals/db"
	"github.com/aasif-10/Olx-API/internals/handlers"
)

func main() {

	cfg := config.MustLoad()

	db, err := db.Connect(cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("main.db.connect %v", err)
	}

	fmt.Println("database connected")
	fmt.Println("starting olx server...")

	lh := handlers.NewListingHandler(db)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", handlers.Health)
	mux.HandleFunc("GET /listings", lh.Listings)
	mux.HandleFunc("DELETE /listings/{id}", lh.Delete)

	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	log.Printf("server is listening on %s", srv.Addr)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Server error")
	}
}
