package api

import (
	"net/http"
)

type api struct {
	captureInProgress bool
	captureErr        error
}

// Start creates an instance of the api, and serves it
func Start() error {
	a := &api{}

	return http.ListenAndServe(":3333", a.createRouter())
}
