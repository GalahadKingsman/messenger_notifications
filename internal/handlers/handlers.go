package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/GalahadKingsman/messenger_notifications/internal/auth"
	"github.com/GalahadKingsman/messenger_notifications/subscriber"
	"log"
	"net/http"
	"strings"
	"time"
)

func LongPollHandler(w http.ResponseWriter, r *http.Request) {

	token := extractToken(r.Header.Get("Authorization"))
	if token == "" {
		http.Error(w, "missing or invalid token", http.StatusUnauthorized)
		return
	}

	userID, err := auth.ExtractUserID(token)
	if err != nil {
		log.Printf("[Notifications] auth failed: %v", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("[Notifications] LongPollHandler hit for userID=%s, errOnExtract=%v", userID, err)

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	notifs, err := subscriber.WaitForMessages(ctx, userID)
	if err != nil {

		if errors.Is(err, context.DeadlineExceeded) {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		log.Printf("[Notifications] WaitForMessages error: %v", err)
		http.Error(w, "gateway error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifs)
}

func extractToken(header string) string {
	if !strings.HasPrefix(header, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(header, "Bearer ")
}
