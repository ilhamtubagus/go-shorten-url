package routes

import (
	"fmt"
	"github.com/ilhamtubagus/go-shorten-url/services"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Routes struct {
	template *template.Template
	service  services.ShortenedService
}

func NewRoutes(t *template.Template, s services.ShortenedService) *Routes {
	return &Routes{template: t, service: s}
}

func (routes *Routes) Index() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		err := routes.template.ExecuteTemplate(w, "index.html", nil)

		if err != nil {
			log.Print(err)
		}
	}
}

func (routes *Routes) NotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := routes.template.ExecuteTemplate(w, "404.html", nil)

		if err != nil {
			log.Print(err)
		}
	}
}

func (routes *Routes) ShortenURL() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		originalURL := r.FormValue("originalURL")
		shortenedURL, err := routes.service.ShortenURL(originalURL)

		host := fmt.Sprintf("%s:%s", os.Getenv("SERVICE_HOST"), os.Getenv("SERVICE_PORT"))
		shortenedURL.ShortenedURL = fmt.Sprintf("%s/%s", host, shortenedURL.ShortenedURL)

		if err != nil {
			log.Print(err)
		}

		err = routes.template.ExecuteTemplate(w, "shorten.html", shortenedURL)

		if err != nil {
			log.Print(err)
		}
	}
}

func (routes *Routes) RedirectURL() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		shortCode := ps.ByName("shortCode")

		// Get the original URL from the service
		shortenedURL, err := routes.service.GetByShortCode(shortCode)
		if err != nil {
			err := routes.template.ExecuteTemplate(w, "404.html", nil)

			if err != nil {
				log.Print(err)
			}

			return
		}

		// Redirect to the original URL
		http.Redirect(w, r, shortenedURL.OriginalURL, http.StatusFound)
	}
}
