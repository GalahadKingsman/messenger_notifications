package redis

import (
	"context"
	"os"

	redis "github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func NewClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"), // например, "localhost:6379"
		Password: "",
		DB:       0,
	})
}

func AddNotification(rdb *redis.Client, userID string, message string) error {
	return rdb.LPush(Ctx, "notifications:"+userID, message).Err()
}

func GetNotifications(rdb *redis.Client, userID string) ([]string, error) {
	return rdb.LRange(Ctx, "notifications:"+userID, 0, -1).Result()
}

func ClearNotifications(rdb *redis.Client, userID string) error {
	return rdb.Del(Ctx, "notifications:"+userID).Err()
}

func HasNotifications(rdb *redis.Client, userID string) (bool, error) {
	count, err := rdb.LLen(Ctx, "notifications:"+userID).Result()
	return count > 0, err
}