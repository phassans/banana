package controller

import (
	"context"
	"net/url"
)

type listingsEndpoint struct{}

func (r listingsEndpoint) GetPath() string { return "/listings" }
func (r listingsEndpoint) Do(ctx context.Context, rtr *router, values url.Values) (interface{}, error) {
	permissions, err := getListings()
	return permissions, err
}

func getListings() ([]string, error) {
	listings := []string{"listing1"}
	return listings, nil
}

var allListings getEndPoint = listingsEndpoint{}
