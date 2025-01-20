package main

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
)

var routes = []string{
	// populate this slice with the routes from the routes.go file
	"/",
	"/login",
	"/logout",
	"/register",
	"/activate-account",
	"/members/plans",
	"/members/subscribe",
}

func Test_Routes_Exist(t *testing.T) {
	// call the routes method which returns a http.Handler
	testRoutes := testApp.routes()
	chiRoutes := testRoutes.(chi.Router)

	// loop through the routes slice
	for _, route := range routes {
		routeExists(t, chiRoutes, route)
	}
}

func routeExists(t *testing.T, routes chi.Router, route string) {
	found := false

	chi.Walk(routes, func(method string, foundRoute string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if foundRoute == route {
			found = true
		}
		return nil
	})

	if !found {
		t.Errorf("route '%s' not found in registered routes", route)
	}
}
