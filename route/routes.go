package route

import (
	"fmt"
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/go-chi/chi"
	"github.com/pshassans/banana/controller"
	"github.com/pshassans/banana/model"
)

func APIServerHandler(engines model.Engine) http.Handler {
	r := newAPIRouter(engines)
	return gziphandler.GzipHandler(r)
}

func newAPIRouter(engines model.Engine) chi.Router {
	r := chi.NewRouter()

	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "application is healthy")
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	r.Mount("/", controller.NewRESTRouter(engines))

	return r
}
