package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractToken_Valid(t *testing.T) {
	header := "Bearer abc.def.ghi"
	token := extractToken(header)
	assert.Equal(t, "abc.def.ghi", token)
}

func TestExtractToken_EmptyToken(t *testing.T) {
	header := "Bearer "
	token := extractToken(header)
	assert.Equal(t, "", token)
}

func TestExtractToken_InvalidPrefix(t *testing.T) {
	header := "Token abc.def.ghi"
	token := extractToken(header)
	assert.Equal(t, "", token)
}

func TestExtractToken_MissingHeader(t *testing.T) {
	header := ""
	token := extractToken(header)
	assert.Equal(t, "", token)
}
