package slogger

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDevNullLogger(t *testing.T) {
	t.Run("it should do nothing", func(t *testing.T) {
		logger := NewDevNullLogger()
		logger.Info("hello world")
		assert.True(t, true)
	})
}
