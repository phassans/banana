package controller

import (
	"context"
	"net/http"
	"net/url"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-chi/chi"
	"github.com/phassans/banana/helper"
	"github.com/phassans/banana/model"
	"github.com/rs/cors"
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
		businessAll,
		listingInfo,
		businessInfo,
		userGet,
		getStats,
		listingUpdateDate,
	}

	// createEndpoints lists POST endpoints that create records.
	addEndpoints = []postEndpoint{
		userAdd,
		userVerify,
		userEdit,

		businessAdd,
		businessEdit,
		businessDelete,

		listingAdd,
		listingsSearch,
		listingDelete,
		listingAll,
		listingEdit,

		notificationAdd,
		notificationDelete,
		notificationAll,

		favouriteAdd,
		favouriteDelete,
		favouriteAll,

		registerPhone,

		preferenceAdd,
		preferenceDelete,
		preferenceAll,

		contactUs,

		upvoteAdd,
		upvoteDelete,
		upvoteCount,

		happyHourSubmit,
		reportInaccurate,

		fetchUser,
		fetchUserChannels,

		webhook,
	}
)

// NewRESTRouter construct a Router interface for Restful API.
func NewRESTRouter(engines model.Engine) http.Handler {
	rtr := &router{
		engines,
		chi.NewRouter(),
	}

	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	rtr.Use(cors.Handler)

	rtr.Use(
		helper.SetJSONContentResponse,
		helper.SetFieldsInLogger,
	)

	rtr.Route(apiVersion, func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Get("/openapp/{devicetype}/{id}", rtr.LinkHandler())
		})

		for _, endpoint := range getEndpoints {
			r.Group(func(r chi.Router) {
				r.Get(endpoint.GetPath(), rtr.newGetHandler(endpoint))
			})
		}

		for _, endpoint := range addEndpoints {
			if endpoint.GetPath() == "/hhsubmit" {
				r.Group(func(r chi.Router) {
					r.Post(endpoint.GetPath(), rtr.newImageHandler(endpoint))
				})
				continue
			}

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
	return hystrix.Do(endpoint.GetPath(), f, nil)
}
