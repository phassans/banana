package routes

import (
	"fmt"
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/go-chi/chi"
	"github.com/pshassans/banana/api"
	"github.com/pshassans/banana/engine"
)

func APIServerHandler(engines engine.Engine) http.Handler {
	r := newAPIRouter(engines)
	return gziphandler.GzipHandler(r)
}

func newAPIRouter(engines engine.Engine) chi.Router {
	r := chi.NewRouter()

	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "application is healthy")
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	r.Mount("/", api.NewRESTRouter(engines))

	return r
}
