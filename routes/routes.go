package routes

import (
	"fmt"
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/go-chi/chi"
	"github.com/pshassans/banana/api"
	"github.com/pshassans/banana/engine"
)

func APIServerHandler(businessEngine engine.BusinessEngine, ownerEngine engine.OwnerEngine, listingEngine engine.ListingEngine) http.Handler {
	r := newAPIRouter(businessEngine, ownerEngine, listingEngine)
	return gziphandler.GzipHandler(r)
}

func newAPIRouter(businessEngine engine.BusinessEngine, ownerEngine engine.OwnerEngine, listingEngine engine.ListingEngine) chi.Router {
	r := chi.NewRouter()

	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "application is healthy")
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	r.Mount("/", api.NewRESTRouter(businessEngine, ownerEngine, listingEngine))

	return r
}
