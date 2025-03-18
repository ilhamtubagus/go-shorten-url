package services

import (
	"github.com/ilhamtubagus/go-shorten-url/entity"
	"github.com/ilhamtubagus/go-shorten-url/repository"
)

type ShortenedService interface {
	ShortenURL(originalURL string) (*entity.ShortenedURL, error)
	GetByShortCode(shortcode string) (*entity.ShortenedURL, error)
}

type ShortenedServiceIml struct {
	repository repository.ShortenedRepository
}

func NewShortenedService(repo repository.ShortenedRepository) ShortenedService {
	return &ShortenedServiceIml{repository: repo}
}

func (s *ShortenedServiceIml) ShortenURL(originalURL string) (*entity.ShortenedURL, error) {
	url := &entity.ShortenedURL{
		OriginalURL: originalURL,
	}
	url.GenerateShortCode()
	// todo: handle duplicate shortcode

	err := s.repository.Insert(*url)
	if err != nil {
		return nil, err
	}

	return url, nil
}

func (s *ShortenedServiceIml) GetByShortCode(shortcode string) (*entity.ShortenedURL, error) {
	return s.repository.GetByShortCode(shortcode)
}
