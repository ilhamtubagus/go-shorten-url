package main

import (
	"github.com/ilhamtubagus/go-shorten-url/routes"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"log"
	"net/http"
)

func main() {
	tmpl := template.Must(template.New("").ParseGlob("./templates/*"))

	router := httprouter.New()
	routesDefs := routes.NewRoutes(tmpl)

	router.NotFound = http.HandlerFunc(routesDefs.NotFound())

	router.GET("/", routesDefs.Index())
	router.POST("/shorten", routesDefs.ShortenURL())

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Server running on port 8080")
}
