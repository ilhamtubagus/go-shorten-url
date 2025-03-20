package repository

import (
	"context"
	"errors"
	"github.com/ilhamtubagus/go-shorten-url/constants"
	"github.com/ilhamtubagus/go-shorten-url/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"log"
	"os"
	"strconv"
	"time"
)

type ShortenedRepository interface {
	GetByShortCode(ctx context.Context, shortcode string) (*entity.ShortenedURL, error)
	Insert(ctx context.Context, payload entity.ShortenedURL) error
	GetShortenedURLs(ctx context.Context) (*[]entity.ShortenedURL, error)
}

type ShortenedRepositoryIml struct {
	cache      Cache[entity.ShortenedURL]
	col        *mongo.Collection
	cacheTasks chan entity.ShortenedURL
}

func NewShortenedRepository(cache Cache[entity.ShortenedURL], col *mongo.Collection) *ShortenedRepositoryIml {
	repo := &ShortenedRepositoryIml{
		cache:      cache,
		col:        col,
		cacheTasks: make(chan entity.ShortenedURL, 100),
	}

	// Start a fixed number of workers
	for i := 0; i < 5; i++ { // Number of workers, adjust as needed
		go repo.cacheWorker()
	}

	return repo
}

func (i *ShortenedRepositoryIml) cacheWorker() {
	for shortenedURL := range i.cacheTasks {
		i.insertCache(shortenedURL)
	}
}

func (i *ShortenedRepositoryIml) insertCache(shortenedURL entity.ShortenedURL) {
	log.Printf("inserting cache %v", shortenedURL.ShortCode)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ttlString := os.Getenv("REDIS_TTL")
	ttl, err := strconv.Atoi(ttlString)
	if err != nil {
		log.Printf("error parsing int %v\n", err)
	}

	err = i.cache.Put(ctx, shortenedURL.ShortCode, shortenedURL, uint64(ttl))
	if err != nil {
		log.Printf("error inserting cache %v\n", err)
	}

	log.Printf("insert cache %v success \n", shortenedURL.ShortCode)
}

func (i *ShortenedRepositoryIml) GetByShortCode(ctx context.Context, shortcode string) (*entity.ShortenedURL, error) {
	log.Printf("getting from cache %v\n", shortcode)

	shortenedURL, err := i.cache.Get(ctx, shortcode)

	if err != nil {
		if !errors.Is(err, constants.ErrorCacheNotFound) {
			return &entity.ShortenedURL{}, err
		}

		if errors.Is(err, constants.ErrorCacheNotFound) {
			log.Printf("getting from mongodb %v\n", shortcode)

			filter := bson.D{{"shortCode", shortcode}}
			var shortened entity.ShortenedURL
			err := i.col.FindOne(ctx, filter).Decode(&shortened)

			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					return nil, constants.ErrorNotFound
				}

				return nil, err
			}

			i.cacheTasks <- shortened

			return &shortened, nil
		}
	}

	return &shortenedURL, nil
}

func (i *ShortenedRepositoryIml) Insert(ctx context.Context, shortenedURL entity.ShortenedURL) error {
	log.Printf("inserting shortened URL into mongodb %v\n", shortenedURL.ShortCode)

	_, err := i.col.InsertOne(ctx, shortenedURL)
	if err != nil {
		return err
	}

	log.Printf("success insert into database %v\n", shortenedURL.ShortCode)

	i.cacheTasks <- shortenedURL

	return nil
}

func (i *ShortenedRepositoryIml) GetShortenedURLs(ctx context.Context) (*[]entity.ShortenedURL, error) {
	// todo : retrieve from MongoDB
	return nil, nil
}
