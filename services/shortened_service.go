package services

import (
	"context"
	"errors"
	"github.com/ilhamtubagus/go-shorten-url/entity"
	"github.com/ilhamtubagus/go-shorten-url/repository"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"log"
	"strconv"
)

type ShortenedService interface {
	ShortenURL(ctx context.Context, originalURL string) (*entity.ShortenedURL, error)
	GetByShortCode(ctx context.Context, shortcode string) (*entity.ShortenedURL, error)
	ListShortenedURLs(ctx context.Context) (*[]entity.ShortenedURL, error)
	DeleteShortenedURL(ctx context.Context, shortcode string) error
}

type ShortenedServiceIml struct {
	repository repository.ShortenedRepository
}

func NewShortenedService(repo repository.ShortenedRepository) ShortenedService {
	return &ShortenedServiceIml{repository: repo}
}

func (s *ShortenedServiceIml) insertWithRetry(ctx context.Context, originalURL string, attempt int) (*entity.ShortenedURL, error) {
	if attempt > 10 {
		return nil, errors.New("too many duplicate attempts")
	}

	shortened := entity.ShortenedURL{
		OriginalURL: originalURL,
	}
	shortened.GenerateShortCode(strconv.Itoa(attempt))

	err := s.repository.Insert(ctx, shortened)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Printf("attempt %d: duplicate shortCode '%s', retrying...\n", attempt, shortened.ShortCode)

			return s.insertWithRetry(ctx, originalURL, attempt+1)
		}
		return nil, err
	}

	return &shortened, nil
}

func (s *ShortenedServiceIml) ShortenURL(ctx context.Context, originalURL string) (*entity.ShortenedURL, error) {
	shorten, err := s.insertWithRetry(ctx, originalURL, 1)

	if err != nil {
		return nil, err
	}

	_ = shorten.GenerateShortenedURL()

	return shorten, nil
}

func (s *ShortenedServiceIml) GetByShortCode(ctx context.Context, shortcode string) (*entity.ShortenedURL, error) {
	shorten, err := s.repository.GetByShortCode(ctx, shortcode)
	if err != nil {
		return nil, err
	}
	// shortened URL with server host generated on the fly
	// in database we are not saving server host, instead we are only saving short code
	err = shorten.GenerateShortenedURL()
	if err != nil {
		return nil, err
	}

	return shorten, nil
}

func (s *ShortenedServiceIml) ListShortenedURLs(ctx context.Context) (*[]entity.ShortenedURL, error) {
	shortenedURLs, err := s.repository.GetShortenedURLs(ctx)
	if err != nil {
		return nil, err
	}

	for i := range *shortenedURLs {
		_ = (*shortenedURLs)[i].GenerateShortenedURL()
	}

	return shortenedURLs, nil
}

func (s *ShortenedServiceIml) DeleteShortenedURL(ctx context.Context, shortcode string) error {
	err := s.repository.DeleteByShortCode(ctx, shortcode)
	if err != nil {
		return err
	}

	return nil
}
