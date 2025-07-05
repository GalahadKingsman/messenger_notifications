package main

import (
	"github.com/GalahadKingsman/messenger_notifications/internal/handlers"
	"github.com/GalahadKingsman/messenger_notifications/subcriber"
	"log"
	"net/http"
	"os"
)

func main() {
	subcriber.InitRedis()

	http.HandleFunc("/notifications/longpoll", handlers.LongPollHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("[notifications] Listening on :%s...", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
