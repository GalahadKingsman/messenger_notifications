package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	redisLib "github.com/redis/go-redis/v9"
	"messenger_notification/internal/redis"
)

type Handler struct {
	rdb *redisLib.Client
}

func NewHandler(rdb *redisLib.Client) *Handler {
	return &Handler{rdb: rdb}
}

func (h *Handler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	notifs, err := redis.GetNotifications(h.rdb, userID)
	if err != nil {
		http.Error(w, "Failed to get notifications", 500)
		return
	}
	json.NewEncoder(w).Encode(notifs)
}

func (h *Handler) ClearNotifications(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if err := redis.ClearNotifications(h.rdb, userID); err != nil {
		http.Error(w, "Failed to clear notifications", 500)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) LongPollNotifications(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			w.WriteHeader(http.StatusNoContent)
			return
		case <-ticker.C:
			has, err := redis.HasNotifications(h.rdb, userID)
			if err != nil {
				http.Error(w, "Error checking notifications", 500)
				return
			}
			if has {
				notifs, _ := redis.GetNotifications(h.rdb, userID)
				json.NewEncoder(w).Encode(notifs)
				return
			}
		}
	}
}
