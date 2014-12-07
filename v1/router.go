package main

import (
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, r := range routes {
		h := Logger(r.Handler, r.Name)
		router.Methods(r.Method).Path(r.Pattern).Name(r.Name).Handler(h)
	}
	return router
}
