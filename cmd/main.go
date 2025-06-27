package main

import (
	"log"
	"messenger_notification/internal/handlers"
	"messenger_notification/internal/redis"
	"messenger_notification/internal/subscriber"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found, using system environment")
	}

	rdb := redis.NewClient()
	subscriber.Subscribe(rdb)

	r := chi.NewRouter()
	h := handlers.NewHandler(rdb)

	r.Get("/notifications/{userID}", h.GetNotifications)
	r.Get("/notifications/{userID}/longpoll", h.LongPollNotifications)
	r.Delete("/notifications/{userID}", h.ClearNotifications)

	addr := ":8082"
	log.Println("Notifications service listening on", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
