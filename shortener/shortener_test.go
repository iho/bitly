package shortener_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/iho/bitly/shortener/mock_shortener"
	"github.com/stretchr/testify/assert"
)

func TestShortener(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	shortener := mock_shortener.NewMockShortener(ctrl)

	defer ctrl.Finish()

	t.Run("Save", func(t *testing.T) {
		shortener.EXPECT().Save(ctx, "http://test.com").Return("test", nil)
		url, err := shortener.Save(ctx, "http://test.com")
		assert.NoError(t, err)
		assert.Equal(t, "test", url)
	})

	t.Run("Load", func(t *testing.T) {
		shortener.EXPECT().Load(ctx, "test").Return("http://test.com", nil)
		url, err := shortener.Load(ctx, "test")
		assert.NoError(t, err)
		assert.Equal(t, "http://test.com", url)
	})

}
