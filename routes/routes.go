package routes

import (
	"fmt"
	"net/http"

	"context"

	"github.com/NYTimes/gziphandler"
	"github.com/banana/api"
	"github.com/banana/defaults"
	"github.com/banana/helper"
	"github.com/go-chi/chi"
)

func WebServerRouter(ctx context.Context) http.Handler {
	r := newRouter(ctx)
	return gziphandler.GzipHandler(r)
}

func newRouter(ctx context.Context) chi.Router {
	r := chi.NewRouter()

	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "application is healthy")
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	r.Mount("/", api.NewAPIRouter(helper.WithValue(ctx, helper.ApiVersion, defaults.CurrentApiVersion)))

	return r
}
