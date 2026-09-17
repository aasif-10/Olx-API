package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/aasif-10/Olx-API/internals/config"
	"github.com/aasif-10/Olx-API/internals/db"
	"github.com/aasif-10/Olx-API/internals/handlers"
	"github.com/aasif-10/Olx-API/internals/middlewares"
)

func main() {

	cfg := config.MustLoad()

	db, err := db.Connect(cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("main.db.connect %v", err)
	}

	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	fmt.Println("database connected")
	fmt.Println("starting olx server...")

	ah := handlers.NewAuthHandler(db, logger, cfg)
	lh := handlers.NewListingHandler(db, logger)
	requireAuth := middlewares.RequireAuth(logger, cfg.JwtKey)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", handlers.Health)
	mux.HandleFunc("GET /listings", lh.Listings)
	mux.Handle("DELETE /listings/{id}", requireAuth(http.HandlerFunc(lh.Delete)))
	mux.HandleFunc("POST /listings", lh.Create)
	mux.HandleFunc("POST /signup", ah.Signup)
	mux.HandleFunc("POST /signin", ah.Signin)

	handler := middlewares.RequestId(mux)

	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
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
