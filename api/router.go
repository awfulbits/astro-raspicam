package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func (a *api) createRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/camera/{imageSetConfigID}", func(r chi.Router) {
		r.With(a.loadImageSetConfig).Post("/capture", a.capture) // POST /camera/imageSetConfigID/capture
	})

	return r
}

type errResponse struct {
	HTTPStatusCode int    `json:"-"`               // http response status code
	StatusText     string `json:"status"`          // client-level status message
	ErrorText      string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *errResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func errInternalServer(err error) *errResponse {
	return &errResponse{
		HTTPStatusCode: 500,
		StatusText:     "Internal Server Error.",
		ErrorText:      err.Error(),
	}
}

func errNotFound() *errResponse {
	return &errResponse{
		HTTPStatusCode: 404,
		StatusText:     "Resource not found.",
	}
}
