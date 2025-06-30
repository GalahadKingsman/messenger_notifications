package handlers

import (
	"encoding/json"
	"messenger_notification/internal/auth"
	"messenger_notification/subcriber"
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
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	notifs, err := subcriber.WaitForMessages(userID, 30*time.Second)
	if err != nil {
		http.Error(w, "timeout", http.StatusGatewayTimeout)
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
