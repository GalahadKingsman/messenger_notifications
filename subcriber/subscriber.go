package subcriber

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"messenger_notification/internal/models"
	"os"
	"time"
)

var ctx = context.Background()
var client *redis.Client

func InitRedis() {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	client = redis.NewClient(&redis.Options{
		Addr: addr,
	})
}

func WaitForMessages(userID string, timeout time.Duration) ([]models.Notification, error) {
	channel := fmt.Sprintf("notifications:%s", userID)
	sub := client.Subscribe(ctx, channel)
	defer sub.Close()

	ch := sub.Channel()
	select {
	case msg := <-ch:
		var notif models.Notification
		if err := json.Unmarshal([]byte(msg.Payload), &notif); err != nil {
			return nil, err
		}
		return []models.Notification{notif}, nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("timeout")
	}
}
