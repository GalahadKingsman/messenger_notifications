package main

import (
	"log"
	"messenger_notification/internal/handlers"
	"messenger_notification/subcriber"
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
