package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/banana/helper"
	"github.com/go-chi/chi"
)

var (
	apiVersion = "/v1"
)

type router struct {
	chi.Router
}

type (
	endpoint interface {
		GetPath() string
	}

	getEndPoint interface {
		endpoint
		Do(context.Context, *router, url.Values) (interface{}, error)
	}

	postEndpoint interface {
		endpoint
		HTTPRequest() interface{}
		HTTPResult() interface{}
		Execute(context.Context, *router, interface{}) (interface{}, error)
	}
)

var (
	// getEndpoints lists all the GET endpoints.
	getEndpoints = []getEndPoint{
		allListings,
	}

	// createEndpoints lists POST endpoints that create records.
	createEndpoints = []postEndpoint{
		createListing,
	}
)

// NewRestAPIRouter construct a Router interface for Restful API.
func NewRESTRouter() http.Handler {
	rtr := &router{
		chi.NewRouter(),
	}

	rtr.Use(
		helper.SetJSONContentResponse,
	)

	rtr.Route(apiVersion, func(r chi.Router) {
		for _, endpoint := range getEndpoints {
			r.Group(func(r chi.Router) {
				r.Get(endpoint.GetPath(), rtr.newGetHandler(endpoint))
			})
		}

		for _, endpoint := range createEndpoints {
			r.Group(func(r chi.Router) {
				r.Post(endpoint.GetPath(), rtr.newPostHandler(endpoint))
			})
		}
	})

	return rtr
}

func (rtr *router) cleanup(e *error, w http.ResponseWriter) {
	err := *e
	if err != nil {
		e := NewAPIError(err)
		e.Send(w)
	}
}

func hystrixName(endpoint endpoint) string {
	return fmt.Sprintf("%s", endpoint.GetPath())
}

func hystrixCall(endpoint endpoint, f func() error) error {
	return hystrix.Do(hystrixName(endpoint), f, nil)
}
