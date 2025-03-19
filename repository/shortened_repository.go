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
	GetByShortCode(shortcode string) (*entity.ShortenedURL, error)
	Insert(payload entity.ShortenedURL) error
	GetShortenedURLs() (*[]entity.ShortenedURL, error)
}

type ShortenedRepositoryIml struct {
	cache Cache[entity.ShortenedURL]
	col   *mongo.Collection
}

func NewShortenedRepository(cache Cache[entity.ShortenedURL], col *mongo.Collection) *ShortenedRepositoryIml {
	return &ShortenedRepositoryIml{cache: cache, col: col}
}

func (i ShortenedRepositoryIml) insertCache(shortenedURL entity.ShortenedURL) {
	log.Println("inserting cache")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ttlString := os.Getenv("REDIS_TTL")
	ttl, err := strconv.Atoi(ttlString)
	if err != nil {
		log.Printf("error parsing int: %v\n", err)
	}

	err = i.cache.Put(ctx, shortenedURL.ShortenedURL, shortenedURL, uint64(ttl))
	if err != nil {
		log.Printf("error inserting cache: %v\n", err)
	}

	log.Printf("insert cache success, key: %v\n", shortenedURL.ShortenedURL)
}

func (i ShortenedRepositoryIml) GetByShortCode(shortcode string) (*entity.ShortenedURL, error) {
	log.Printf("getting from cache %v\n", shortcode)

	shortenedURL, err := i.cache.Get(context.TODO(), shortcode)

	if err != nil {
		if !errors.Is(err, constants.ErrorCacheNotFound) {
			return &entity.ShortenedURL{}, err
		}

		if errors.Is(err, constants.ErrorCacheNotFound) {
			log.Printf("getting from mongodb %v\n", shortcode)

			filter := bson.D{{"shortenedURL", shortcode}}
			var shortenedURL entity.ShortenedURL
			err := i.col.FindOne(context.TODO(), filter).Decode(&shortenedURL)

			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					return nil, constants.ErrorNotFound
				}

				return nil, err
			}
			go i.insertCache(shortenedURL)

			return &shortenedURL, nil
		}
	}

	return &shortenedURL, nil
}

func (i ShortenedRepositoryIml) Insert(shortenedURL entity.ShortenedURL) error {
	log.Printf("inserting shortened URL into mongodb %v\n", shortenedURL.ShortenedURL)

	_, err := i.col.InsertOne(context.TODO(), shortenedURL)
	if err != nil {
		return err
	}

	log.Printf("success insert into database %v\n", shortenedURL.ShortenedURL)

	go i.insertCache(shortenedURL)

	return nil
}

func (i ShortenedRepositoryIml) GetShortenedURLs() (*[]entity.ShortenedURL, error) {
	// todo : retrieve from MongoDB
	return nil, nil
}
