package routes

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"log"
	"net/http"
)

type Routes struct {
	Template *template.Template
}

func NewRoutes(t *template.Template) *Routes {
	return &Routes{Template: t}
}

func (routes *Routes) Index() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		err := routes.Template.ExecuteTemplate(w, "index.html", nil)

		if err != nil {
			log.Print(err)
		}
	}
}

func (routes *Routes) NotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := routes.Template.ExecuteTemplate(w, "404.html", nil)

		if err != nil {
			log.Print(err)
		}
	}
}

func (routes *Routes) ShortenURL() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		originalURL := r.FormValue("originalURL")

		fmt.Println(originalURL)
	}
}
