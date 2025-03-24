package main

import (
	"fmt"
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

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
}

func main() {
	tmpl := template.Must(template.New("").ParseGlob("./templates/*"))

	redisClient := server.ConnectRedisClient()
	mongoClient := server.ConnectMongoClient()

	shortenRedisCache := repository.NewRedisCache[entity.ShortenedURL](redisClient)
	shortenCollection := mongoClient.Database("shorten").Collection("shorten")
	shortenRepository := repository.NewShortenedRepository(shortenRedisCache, shortenCollection)

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
