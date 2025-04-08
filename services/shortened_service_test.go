package services

import (
	"context"
	"errors"
	"github.com/ilhamtubagus/shortenurl/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"testing"
)

// MockShortenedRepository is a mock type for repository.ShortenedRepository
type MockShortenedRepository struct {
	mock.Mock
}

func (m *MockShortenedRepository) Insert(ctx context.Context, shortened entity.ShortenedURL) error {
	args := m.Called(ctx, shortened)
	return args.Error(0)
}

func (m *MockShortenedRepository) GetByShortCode(ctx context.Context, shortcode string) (*entity.ShortenedURL, error) {
	args := m.Called(ctx, shortcode)
	return args.Get(0).(*entity.ShortenedURL), args.Error(1)
}

func (m *MockShortenedRepository) GetShortenedURLs(ctx context.Context) (*[]entity.ShortenedURL, error) {
	args := m.Called(ctx)
	return args.Get(0).(*[]entity.ShortenedURL), args.Error(1)
}

func (m *MockShortenedRepository) DeleteByShortCode(ctx context.Context, shortcode string) error {
	args := m.Called(ctx, shortcode)
	return args.Error(0)
}

func (m *MockShortenedRepository) UpdateByShortCode(ctx context.Context, shortcode string, originalURL string) (*entity.ShortenedURL, error) {
	args := m.Called(ctx, shortcode, originalURL)
	return args.Get(0).(*entity.ShortenedURL), args.Error(1)
}

func createDuplicateKeyError() error {
	writeErr := mongo.WriteException{
		WriteErrors: []mongo.WriteError{
			{Code: 11000, Message: "duplicate key error", Index: 0},
		},
	}
	return writeErr
}

func TestShortenedServiceIml_ShortenURL(t *testing.T) {
	mockRepo := new(MockShortenedRepository)
	service := NewShortenedService(mockRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		originalURL := "https://example.com"
		mockRepo.On("Insert", ctx, mock.AnythingOfType("entity.ShortenedURL")).Return(nil)

		result, err := service.ShortenURL(ctx, originalURL)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, originalURL, result.OriginalURL)
		assert.NotEmpty(t, result.ShortCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("DuplicateKeyError", func(t *testing.T) {
		originalURL := "https://example.com"
		mockRepo.On("Insert", ctx, mock.AnythingOfType("entity.ShortenedURL")).
			Return(createDuplicateKeyError()).
			Times(10)

		result, err := service.ShortenURL(ctx, originalURL)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "too many duplicate attempts")
		mockRepo.AssertExpectations(t)
	})
}

func TestShortenedServiceIml_GetByShortCode(t *testing.T) {
	mockRepo := new(MockShortenedRepository)
	service := NewShortenedService(mockRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		shortcode := "abc123"
		expectedURL := &entity.ShortenedURL{ShortCode: shortcode, OriginalURL: "https://example.com"}
		mockRepo.On("GetByShortCode", ctx, shortcode).Return(expectedURL, nil)

		result, err := service.GetByShortCode(ctx, shortcode)

		assert.NoError(t, err)
		assert.Equal(t, expectedURL, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		shortcode := "notfound"
		mockRepo.On("GetByShortCode", ctx, shortcode).Return((*entity.ShortenedURL)(nil), errors.New("not found"))

		result, err := service.GetByShortCode(ctx, shortcode)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestShortenedServiceIml_ListShortenedURLs(t *testing.T) {
	mockRepo := new(MockShortenedRepository)
	service := NewShortenedService(mockRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		expectedURLs := &[]entity.ShortenedURL{
			{ShortCode: "abc123", OriginalURL: "https://example1.com"},
			{ShortCode: "def456", OriginalURL: "https://example2.com"},
		}
		mockRepo.On("GetShortenedURLs", ctx).Return(expectedURLs, nil)

		result, err := service.ListShortenedURLs(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedURLs, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo.On("GetShortenedURLs", ctx).Return((*[]entity.ShortenedURL)(nil), errors.New("database error"))

		result, err := service.ListShortenedURLs(ctx)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestShortenedServiceIml_DeleteShortenedURL(t *testing.T) {
	mockRepo := new(MockShortenedRepository)
	service := NewShortenedService(mockRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		shortcode := "abc123"
		mockRepo.On("DeleteByShortCode", ctx, shortcode).Return(nil)

		err := service.DeleteShortenedURL(ctx, shortcode)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		shortcode := "notfound"
		mockRepo.On("DeleteByShortCode", ctx, shortcode).Return(errors.New("not found"))

		err := service.DeleteShortenedURL(ctx, shortcode)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestShortenedServiceIml_UpdateShortenedURL(t *testing.T) {
	mockRepo := new(MockShortenedRepository)
	service := NewShortenedService(mockRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		shortcode := "abc123"
		originalURL := "https://newexample.com"
		expectedURL := &entity.ShortenedURL{ShortCode: shortcode, OriginalURL: originalURL}
		mockRepo.On("UpdateByShortCode", ctx, shortcode, originalURL).Return(expectedURL, nil)

		result, err := service.UpdateShortenedURL(ctx, shortcode, originalURL)

		assert.NoError(t, err)
		assert.Equal(t, expectedURL, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		shortcode := "notfound"
		originalURL := "https://newexample.com"
		mockRepo.On("UpdateByShortCode", ctx, shortcode, originalURL).Return((*entity.ShortenedURL)(nil), errors.New("not found"))

		result, err := service.UpdateShortenedURL(ctx, shortcode, originalURL)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}
