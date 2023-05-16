package storage_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/iho/bitly/storage/mock_storage"
	"github.com/stretchr/testify/assert"
)

// generate golang  unit test for storage interface

func TestStorage(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	storage := mock_storage.NewMockStorage(ctrl)

	defer ctrl.Finish()

	t.Run("Save", func(t *testing.T) {
		storage.EXPECT().Save(ctx, "test", "http://test.com").Return(nil)
		err := storage.Save(ctx, "test", "http://test.com")
		assert.NoError(t, err)
	})

	t.Run("Load", func(t *testing.T) {
		storage.EXPECT().Load(ctx, "test").Return("http://test.com", nil)
		url, err := storage.Load(ctx, "test")
		assert.NoError(t, err)
		assert.Equal(t, "http://test.com", url)
	})
}
