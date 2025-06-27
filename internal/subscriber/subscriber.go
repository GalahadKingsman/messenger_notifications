package subscriber

import (
	"encoding/json"
	"log"
	"messenger_notification/internal/redis"

	rds "github.com/redis/go-redis/v9"
)

type Event struct {
	Type    string `json:"type"`
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

func Subscribe(rdb *rds.Client) {
	pubsub := rdb.Subscribe(redis.Ctx, "new_message")
	ch := pubsub.Channel()

	go func() {
		for msg := range ch {
			var event Event
			if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
				log.Println("invalid event:", err)
				continue
			}

			if err := redis.AddNotification(rdb, event.UserID, event.Message); err != nil {
				log.Println("failed to add notification:", err)
			}
		}
	}()
}
