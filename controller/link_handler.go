package controller

import (
	"fmt"
	"net/http"

	"github.com/phassans/banana/shared"

	"github.com/go-chi/chi"
)

const (
	HTML_CONTENT = `<html><head></head><a style="font-size:30px;" href="%s">Open in hungryhour app? click here!</a></html>`
)

func (rtr *router) LinkHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := shared.GetLogger()
		id := chi.URLParam(r, "id")
		deviceType := chi.URLParam(r, "devicetype")

		var link string
		if deviceType == "android" {
			link = fmt.Sprintf("http://hungryhour/%s", id)
		} else if deviceType == "ios" {
			link = fmt.Sprintf("hungryhour://%s", id)
		}

		logger.Info().Msgf("app %s, id %s and URL: %s", deviceType, id, link)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		h := fmt.Sprintf(HTML_CONTENT, link)
		fmt.Fprintln(w, h)
	}
}
