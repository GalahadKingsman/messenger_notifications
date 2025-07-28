package auth

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func generateToken(secret string, claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte(secret))
	return tokenStr
}

func TestExtractUserID_String(t *testing.T) {
	os.Setenv("JWT_SECRET", "mysecret")
	token := generateToken("mysecret", jwt.MapClaims{
		"user_id": "12345",
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	})

	uid, err := ExtractUserID(token)
	assert.NoError(t, err)
	assert.Equal(t, "12345", uid)
}

func TestExtractUserID_Number(t *testing.T) {
	os.Setenv("JWT_SECRET", "mysecret")
	token := generateToken("mysecret", jwt.MapClaims{
		"user_id": 67890,
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	})

	uid, err := ExtractUserID(token)
	assert.NoError(t, err)
	assert.Equal(t, "67890", uid)
}

func TestExtractUserID_MissingUserID(t *testing.T) {
	os.Setenv("JWT_SECRET", "mysecret")
	token := generateToken("mysecret", jwt.MapClaims{
		"exp": time.Now().Add(1 * time.Hour).Unix(),
	})

	uid, err := ExtractUserID(token)
	assert.Error(t, err)
	assert.Equal(t, "", uid)
}

func TestExtractUserID_InvalidToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "mysecret")
	uid, err := ExtractUserID("not.a.token")
	assert.Error(t, err)
	assert.Equal(t, "", uid)
}

func TestExtractUserID_WrongSignature(t *testing.T) {
	os.Setenv("JWT_SECRET", "correct")
	token := generateToken("wrong", jwt.MapClaims{
		"user_id": "1",
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	})

	uid, err := ExtractUserID(token)
	assert.Error(t, err)
	assert.Equal(t, "", uid)
}

func TestExtractUserID_NoSecret(t *testing.T) {
	os.Unsetenv("JWT_SECRET")
	token := generateToken("anything", jwt.MapClaims{
		"user_id": "1",
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	})

	uid, err := ExtractUserID(token)
	assert.EqualError(t, err, "JWT_SECRET not set")
	assert.Equal(t, "", uid)
}
