package routes

import (
	"fmt"
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/banana/api"
	"github.com/go-chi/chi"
)

func APIServerHandler() http.Handler {
	r := newAPIRouter()
	return gziphandler.GzipHandler(r)
}

func newAPIRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "application is healthy")
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	r.Mount("/", api.NewRESTRouter())

	return r
}
