package controller

import (
	"fmt"
	"net/http"

	"github.com/phassans/banana/shared"

	"github.com/go-chi/chi"
)

func (rtr *router) RedirectHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := shared.GetLogger()
		id := chi.URLParam(r, "id")
		redirectURL := fmt.Sprintf("http://hungryhour/%s", id)
		logger.Info().Msgf("redirecting to URL: %s", redirectURL)
		http.Redirect(w, r, redirectURL, 301)
	}
}
