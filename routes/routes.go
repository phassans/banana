package routes

import (
	"fmt"
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/banana/api"
	"github.com/banana/engine"
	"github.com/go-chi/chi"
)

func APIServerHandler(engine engine.ListingEngine) http.Handler {
	r := newAPIRouter(engine)
	return gziphandler.GzipHandler(r)
}

func newAPIRouter(engine engine.ListingEngine) chi.Router {
	r := chi.NewRouter()

	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "application is healthy")
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	r.Mount("/", api.NewRESTRouter(engine))

	return r
}
