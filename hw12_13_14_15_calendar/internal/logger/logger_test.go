package logger

import (
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestLogger(t *testing.T) {
	var z zapcore.Level

	t.Run("NewLogger create", func(t *testing.T) {
		l, err := NewLogger(z, "/dev/null")
		require.NoError(t, err)
		require.NotNil(t, l)
	})
}
