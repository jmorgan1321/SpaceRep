package main

// import (
// 	"net/http"
// )

// type Route struct {
// 	Name, Method, Pattern string
// 	Handler               http.HandlerFunc
// }

// type Routes []Routes

// const (
// 	apiPrefix = "/api/v1"
// )

// var routes = []Route{
// 	Route{
// 		Name:    "Index",
// 		Method:  "GET",
// 		Pattern: "/",
// 		Handler: nil,
// 	},
// 	Route{
// 		Name:    "CardIndex",
// 		Method:  "GET",
// 		Pattern: apiPrefix + "/cards",
// 		Handler: cardIndexHandler,
// 	},
// 	Route{
// 		Name:    "Review",
// 		Method:  "GET",
// 		Pattern: apiPrefix + "/review/{value}",
// 		Handler: reviewHandler,
// 	},
// 	Route{
// 		Name:    "Save",
// 		Method:  "GET",
// 		Pattern: apiPrefix + "/save",
// 		Handler: saveHandler,
// 	},
// }
