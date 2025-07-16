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
	http.HandleFunc("/notifications", subcriber.PostNotificationHandler)

	port := os.Getenv("PORT")

	log.Println("[notifications] Listening")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
