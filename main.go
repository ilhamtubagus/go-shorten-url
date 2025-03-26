package main

import (
	"fmt"
	goenv "github.com/ilhamtubagus/go-env"
	"github.com/ilhamtubagus/go-shorten-url/config"
	"github.com/ilhamtubagus/go-shorten-url/entity"
	"github.com/ilhamtubagus/go-shorten-url/repository"
	"github.com/ilhamtubagus/go-shorten-url/routes"
	"github.com/ilhamtubagus/go-shorten-url/server"
	"github.com/ilhamtubagus/go-shorten-url/services"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"log"
	"net/http"
	"os"
)

var appConfig = &config.Config{}

func init() {
	if os.Getenv("ENV") == "development" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("error loading .env file")
		}

		log.Println("successfully loaded development environment variables")
	}

	err := goenv.Unmarshal(appConfig)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	tmpl := template.Must(template.New("").ParseGlob("./templates/*"))

	redisClient := server.ConnectRedisClient(appConfig.Redis)
	mongoClient := server.ConnectMongoClient(appConfig.Mongo)

	shortenRedisCache := repository.NewRedisCache[entity.ShortenedURL](redisClient)
	shortenCollection := mongoClient.Database("shorten").Collection("shorten")
	shortenRepository := repository.NewShortenedRepository(shortenRedisCache, shortenCollection, *appConfig)

	shortenService := services.NewShortenedService(shortenRepository)

	router := httprouter.New()
	routesDefs := routes.NewRoutes(tmpl, shortenService)

	router.NotFound = http.HandlerFunc(routesDefs.NotFound())

	router.GET("/", routesDefs.Index())
	router.POST("/shorten-url", routesDefs.ShortenURL())
	router.GET("/shorten-url", routesDefs.ListShortenedURLs())
	router.DELETE("/:shortCode", routesDefs.DeleteShortenedURL())
	router.PATCH("/:shortCode", routesDefs.UpdateShortenedURL())
	router.GET("/s/:shortCode", routesDefs.RedirectURL())

	host := fmt.Sprintf("%s:%s", os.Getenv("SERVICE_HOST"), os.Getenv("SERVICE_PORT"))

	log.Printf("server running on %s\n", host)

	err := http.ListenAndServe(host, router)
	if err != nil {
		log.Fatal(err)
	}
}
