package main

import (
	"net/http"
)

type Route struct {
	Name, Method, Pattern string
	Handler               http.HandlerFunc
}

type Routes []Routes

const (
	apiPrefix = "/api/v1"
)

var routes = []Route{
	Route{
		Name:    "CardIndex",
		Method:  "GET",
		Pattern: apiPrefix + "/cards",
		Handler: cardIndexHandler,
	},
	Route{
		Name:    "Review",
		Method:  "POST",
		Pattern: apiPrefix + "/review/{value}",
		// TODO: restrict path
		// Pattern: apiPrefix + "/review/{value:^(accept|forgot)$}",
		Handler: reviewHandler,
	},
	Route{
		Name:    "Save",
		Method:  "POST",
		Pattern: apiPrefix + "/save",
		Handler: saveHandler,
	},
}
