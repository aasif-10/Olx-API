package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"
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
	db *sql.DB
}

func NewListingHandler(db *sql.DB) *ListingHandler {
	return &ListingHandler{
		db: db,
	}
}

func (lh ListingHandler) Listings(w http.ResponseWriter, r *http.Request) {

	rows, err := lh.db.Query(`
			SELECT id, title, description, price, city, created_at
			FROM listings
			ORDER BY created_at DESC
			LIMIT 100`)

	if err != nil {
		log.Printf("db.query: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	listings := []listing{}

	for rows.Next() {
		var l listing
		err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.City, &l.CreatedAt)

		if err != nil {
			log.Printf("rows.scan: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		listings = append(listings, l)

	}

	err = rows.Err()
	if err != nil {
		log.Printf("rows.err: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

func (lh ListingHandler) Delete(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	_, err := lh.db.Exec(`
			DELETE FROM listings WHERE id = $1`, id)
	if err != nil {
		log.Printf("db.exec: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
