package subcriber

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/GalahadKingsman/messenger_notifications/internal/models"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"os"
	"time"
)

var client *redis.Client
var subscriber *redis.Client

func InitRedis() {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "redis:6379"
	}
	cmdClient := redis.NewClient(&redis.Options{
		Addr:         addr,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := cmdClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("redis ping failed: %v", err)
	}
	subClient := redis.NewClient(&redis.Options{
		Addr:         addr,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  0, //
		WriteTimeout: 0,
	})
	client = cmdClient
	subscriber = subClient
}
func WaitForMessages(ctx context.Context, userID string) ([]models.Notification, error) {
	channelName := fmt.Sprintf("notifications:%s", userID)

	pubsub := subscriber.Subscribe(ctx, channelName)
	defer pubsub.Close()

	if _, err := pubsub.Receive(ctx); err != nil {
		return nil, fmt.Errorf("failed to subscribe to %s: %w", channelName, err)
	}

	msgCh := pubsub.Channel()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case msg := <-msgCh:
		var notif models.Notification
		if err := json.Unmarshal([]byte(msg.Payload), &notif); err != nil {
			log.Printf("bad payload from Redis: %v", err)
			return nil, err
		}
		return []models.Notification{notif}, nil
	}
}

func PostNotificationHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		http.Error(w, "userID query param required", http.StatusBadRequest)
		return
	}

	var notif models.Notification
	if err := json.NewDecoder(r.Body).Decode(&notif); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	payload, _ := json.Marshal(notif)
	channel := fmt.Sprintf("notifications:%s", userID)
	log.Printf("[Notifications] Publishing to channel=%s payload=%s", channel, string(payload))
	if err := client.Publish(r.Context(), channel, payload).Err(); err != nil {
		log.Printf("[Notifications] Publish failed: %v", err)
		http.Error(w, "publish failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
