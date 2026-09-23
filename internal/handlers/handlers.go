package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	PG  *pgxpool.Pool
	RDB *redis.Client
}

func New(pg *pgxpool.Pool, rdb *redis.Client) *Handler {
	return &Handler{PG: pg, RDB: rdb}
}

// POST/GET /cache/set?key=foo&value=bar
func (h *Handler) CacheSet(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	if key == "" {
		http.Error(w, "key is required", http.StatusBadRequest)
		return
	}
	if err := h.RDB.Set(r.Context(), key, value, 0).Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{"status": "ok", "key": key, "value": value})
}

// GET /cache/get?key=foo
func (h *Handler) CacheGet(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "key is required", http.StatusBadRequest)
		return
	}
	val, err := h.RDB.Get(r.Context(), key).Result()
	if err == redis.Nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{"key": key, "value": val})
}

// GET /users
func (h *Handler) Users(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	rows, err := h.PG.Query(ctx, "SELECT id, name FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}
	writeJSON(w, users)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

// GET /healthz 
func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]string{"status": "ok"})
}

// GET /readyz 
func (h *Handler) Readyz(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := h.PG.Ping(ctx); err != nil {
		http.Error(w, "postgres not ready: "+err.Error(), http.StatusServiceUnavailable)
		return
	}
	if err := h.RDB.Ping(ctx).Err(); err != nil {
		http.Error(w, "redis not ready: "+err.Error(), http.StatusServiceUnavailable)
		return
	}
	writeJSON(w, map[string]string{"status": "ready"})
}