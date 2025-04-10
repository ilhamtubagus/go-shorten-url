package repository

import (
	"context"
	"errors"
	"github.com/ilhamtubagus/shortenurl/config"
	"github.com/ilhamtubagus/shortenurl/constants"
	"github.com/ilhamtubagus/shortenurl/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"time"
)

type ShortenedRepository interface {
	GetByShortCode(ctx context.Context, shortcode string) (*entity.ShortenedURL, error)
	Insert(ctx context.Context, payload entity.ShortenedURL) error
	GetShortenedURLs(ctx context.Context) (*[]entity.ShortenedURL, error)
	DeleteByShortCode(ctx context.Context, shortCode string) error
	UpdateByShortCode(ctx context.Context, shortCode string, newOriginalURL string) (*entity.ShortenedURL, error)
}

type ShortenedRepositoryIml struct {
	cache      Cache[entity.ShortenedURL]
	col        *mongo.Collection
	cacheTasks chan entity.ShortenedURL
	config     config.Config
}

func NewShortenedRepository(cache Cache[entity.ShortenedURL], col *mongo.Collection, config config.Config) *ShortenedRepositoryIml {
	repo := &ShortenedRepositoryIml{
		cache:      cache,
		col:        col,
		cacheTasks: make(chan entity.ShortenedURL, 100),
		config:     config,
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

	err := i.cache.Delete(ctx, shortenedURL.ShortCode)
	if err != nil {
		log.Printf("error deleting cache %v\n", err)
	}

	err = i.cache.Put(ctx, shortenedURL.ShortCode, shortenedURL, uint64(i.config.Redis.TTL))
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
	log.Println("getting all shortened URLs from mongodb")

	var shortenedURLs []entity.ShortenedURL
	filter := bson.D{}
	cursor, err := i.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &shortenedURLs); err != nil {
		panic(err)
	}

	log.Println("success get all shortened URLs from mongodb")

	return &shortenedURLs, nil
}

func (i *ShortenedRepositoryIml) DeleteByShortCode(ctx context.Context, shortCode string) error {
	filter := bson.D{{"shortCode", shortCode}}
	if _, err := i.col.DeleteOne(ctx, filter); err != nil {
		return err
	}

	err := i.cache.Delete(ctx, shortCode)
	if err != nil {
		return err
	}

	return nil
}

func (i *ShortenedRepositoryIml) UpdateByShortCode(ctx context.Context, shortCode string, newOriginalURL string) (*entity.ShortenedURL, error) {
	filter := bson.D{{"shortCode", shortCode}}
	update := bson.D{{"$set", bson.D{{"originalURL", newOriginalURL}}}}
	var shortened entity.ShortenedURL

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	err := i.col.FindOneAndUpdate(ctx, filter, update, opts).Decode(&shortened)
	if err != nil {
		return nil, err
	}

	i.cacheTasks <- shortened

	return &shortened, nil
}
