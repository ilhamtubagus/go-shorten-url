package repository

import (
	"context"
	"github.com/ilhamtubagus/go-shorten-url/entity"
	"os"
	"strconv"
)

type ShortenedRepository interface {
	GetByShortCode(shortcode string) (*entity.ShortenedURL, error)
	Insert(payload entity.ShortenedURL) error
	GetShortenedURLs() (*[]entity.ShortenedURL, error)
}

type ShortenedRepositoryIml struct {
	cache Cache[entity.ShortenedURL]
}

func NewShortenedRepository(cache Cache[entity.ShortenedURL]) *ShortenedRepositoryIml {
	return &ShortenedRepositoryIml{cache: cache}
}

func (i ShortenedRepositoryIml) GetByShortCode(shortcode string) (*entity.ShortenedURL, error) {
	shortenedURL, err := i.cache.Get(context.Background(), shortcode)
	if err != nil {
		return &entity.ShortenedURL{}, err
	}

	return &shortenedURL, nil
}

func (i ShortenedRepositoryIml) Insert(payload entity.ShortenedURL) error {
	ttlString := os.Getenv("REDIS_TTL")
	ttl, err := strconv.Atoi(ttlString)
	if err != nil {
		return err
	}

	err = i.cache.Put(context.Background(), payload.ShortenedURL, payload, uint64(ttl))
	if err != nil {
		return err
	}

	return nil
}

func (i ShortenedRepositoryIml) GetShortenedURLs() (*[]entity.ShortenedURL, error) {
	return nil, nil
}
