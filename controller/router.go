package controller

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-chi/chi"
	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/model"
)

var (
	apiVersion = "/v1"
)

type router struct {
	engines model.Engine
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
		Execute(context.Context, *router, interface{}) (interface{}, error)
		Validate(interface{}) error
	}
)

var (
	// getEndpoints lists all the GET endpoints.
	getEndpoints = []getEndPoint{
		listing,
	}

	// createEndpoints lists POST endpoints that create records.
	addEndpoints = []postEndpoint{
		addUser,
		addBusiness,
		addListing,
		listingsSearch,
		verifyUser,
		allListingAdmin,
		addFavorite,
		deleteFavorite,
		favoritesView,
	}
)

// NewRestAPIRouter construct a Router interface for Restful API.
func NewRESTRouter(engines model.Engine) http.Handler {
	rtr := &router{
		engines,
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

		for _, endpoint := range addEndpoints {
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

func hystrixCall(endpoint endpoint, f func() error) error {
	name := fmt.Sprintf("%s", endpoint.GetPath())
	return hystrix.Do(name, f, nil)
}
