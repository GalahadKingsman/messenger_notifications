package subscriber

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/GalahadKingsman/messenger_notifications/internal/models"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func setupTestRedis(t *testing.T) *miniredis.Miniredis {
	srv, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	client = redis.NewClient(&redis.Options{
		Addr: srv.Addr(),
	})
	return srv
}

func TestPostNotificationHandler_Success(t *testing.T) {
	srv := setupTestRedis(t)
	defer srv.Close()

	ctx := context.Background()

	// подписываемся на канал заранее
	redisClient := redis.NewClient(&redis.Options{Addr: srv.Addr()})
	sub := redisClient.Subscribe(ctx, "notifications:42")
	defer sub.Close()

	// теперь вызываем обработчик, который опубликует туда сообщение
	notif := models.Notification{
		From:     "1",
		Message:  "hello world",
		DialogID: 100,
	}
	body, _ := json.Marshal(notif)

	req := httptest.NewRequest(http.MethodPost, "/?userID=42", bytes.NewReader(body))
	w := httptest.NewRecorder()

	PostNotificationHandler(w, req)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// ждём сообщение из канала
	select {
	case msg := <-sub.Channel():
		assert.Contains(t, msg.Channel, "notifications:42")
		assert.Contains(t, msg.Payload, "hello world")
	case <-time.After(time.Second):
		t.Fatal("no message received from Redis channel")
	}
}

func TestPostNotificationHandler_MissingUserID(t *testing.T) {
	srv := setupTestRedis(t)
	defer srv.Close()

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()

	PostNotificationHandler(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPostNotificationHandler_InvalidJSON(t *testing.T) {
	srv := setupTestRedis(t)
	defer srv.Close()

	req := httptest.NewRequest(http.MethodPost, "/?userID=123", bytes.NewBufferString("bad_json"))
	w := httptest.NewRecorder()

	PostNotificationHandler(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
