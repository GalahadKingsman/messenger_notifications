package subcriber

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/GalahadKingsman/messenger_notifications/internal/models"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

var client *redis.Client

func InitRedis() {
	addr := "redis:6379"
	opts := &redis.Options{
		Addr:         addr,
		ReadTimeout:  35 * time.Second,
		WriteTimeout: 35 * time.Second,
		DialTimeout:  5 * time.Second,
	}
	client = redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("redis ping failed: %v", err)
	}
}
func WaitForMessages(ctx context.Context, userID string) ([]models.Notification, error) {
	channel := fmt.Sprintf("notifications:%s", userID)
	sub := client.Subscribe(ctx, channel)
	defer sub.Close()

	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		return nil, err
	}

	var notif models.Notification
	if err := json.Unmarshal([]byte(msg.Payload), &notif); err != nil {
		log.Printf("bad payload from Redis: %v", err)
		return nil, err
	}
	return []models.Notification{notif}, nil
}
