package shortener

import (
	"context"
	"database/sql"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
	"url-shortener/internal/db"
)

type DBQueries interface {
	CreateURL(ctx context.Context, arg db.CreateURLParams) (sql.Result, error)
	GetURLByCode(ctx context.Context, shortCode string) (db.Url, error)
	RecordClick(ctx context.Context, arg db.RecordClickParams) error
}

type Handler struct {
	DB DBQueries
}

func NewHandler(database DBQueries) *Handler {
	return &Handler{DB: database}
}

type ShortenRequest struct {
	LongULR string `json:"long_url"`
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req ShortenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	code := string(b)

	_, err := h.DB.CreateURL(r.Context(), db.CreateURLParams{
		ShortCode:   code,
		OriginalUrl: req.LongULR,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unable to save link"})
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"short_code": "http://127.0.0.1:8080/r/" + code,
		"code":       code,
	})
}
