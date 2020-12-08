package memorystorage

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStorage(t *testing.T) {

	t.Run("Invalid credentials error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		dsn := "host=localhost port=5432 user=user password=password dbname=exampledb sslmode=disable"
		storage := new(Storage)
		err := storage.Connect(ctx, dsn)
		require.Error(t, err)

	})
}
