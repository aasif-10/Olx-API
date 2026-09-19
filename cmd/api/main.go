package main

import (
	"context"
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
	"github.com/aasif-10/Olx-API/internals/storage"
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

	store, err := storage.NewSupabase(context.TODO(), storage.SupabaseConfig{
		ProjectId:    cfg.StorageProjectId,
		AccessKey:    cfg.StorageAccessKey,
		AccessSecret: cfg.StorageAccessSecret,
		Bucket:       cfg.StorageBucket,
	})
	if err != nil {
		log.Fatalf("main.storage.supabase %v", err)
	}

	fmt.Println("storage initiliased...")
	fmt.Println("database connected")
	fmt.Println("starting olx server...")

	ah := handlers.NewAuthHandler(db, logger, cfg)
	lh := handlers.NewListingHandler(db, logger)
	uh := handlers.NewUploadHandler(logger, store)
	requireAuth := middlewares.RequireAuth(logger, cfg.JwtKey)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", handlers.Health)
	mux.HandleFunc("GET /listings", lh.Listings)
	mux.Handle("DELETE /listings/{id}", requireAuth(http.HandlerFunc(lh.Delete)))
	mux.Handle("POST /listings", requireAuth(http.HandlerFunc(lh.Create)))
	mux.HandleFunc("POST /signup", ah.Signup)
	mux.HandleFunc("POST /signin", ah.Signin)
	mux.Handle("POST /upload/presign", requireAuth(http.HandlerFunc(uh.Presign)))

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
