package slogger

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrorAttr(t *testing.T) {
	t.Run("it should return empty if no error given", func(t *testing.T) {
		got := ErrorAttr(nil)
		assert.Equal(t, "error", got.Key)
		assert.Equal(t, "", got.Value.String())
	})
	t.Run("it should the error in a log format", func(t *testing.T) {
		got := ErrorAttr(assert.AnError)
		assert.Equal(t, "error", got.Key)
		assert.Equal(t, assert.AnError.Error(), got.Value.String())
	})
}
