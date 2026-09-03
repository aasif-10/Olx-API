package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/aasif-10/Olx-API/internals/middlewares"
)

type listing struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       string    `json:"price"`
	City        string    `json:"city"`
	CreatedAt   time.Time `json:"created_at"`
}

type ListingHandler struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewListingHandler(db *sql.DB, logger *slog.Logger) *ListingHandler {
	return &ListingHandler{
		db:     db,
		logger: logger,
	}
}

func (lh ListingHandler) Listings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rows, err := lh.db.QueryContext(ctx, `
			SELECT id, title, description, price, city, created_at
			FROM listings
			ORDER BY created_at DESC
			LIMIT 100`)

	if err != nil {
		lh.logger.Error("listings query error", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	listings := []listing{}

	for rows.Next() {
		var l listing
		err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.City, &l.CreatedAt)

		if err != nil {
			lh.logger.Error("row scan error", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		lh.logger.Info("listings fetched", "total", len(listings))
		listings = append(listings, l)

	}

	err = rows.Err()
	if err != nil {
		lh.logger.Error("rows error", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(listings)
}

func (lh ListingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middlewares.RequestIDFromContext(ctx)
	id := r.PathValue("id")

	_, err := lh.db.ExecContext(ctx, `
			DELETE FROM listing WHERE id = $1`, id)
	if err != nil {
		lh.logger.Error("delete", "listings", id, "request_id", requestId, "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
