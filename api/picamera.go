package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/awfulbits/astro-raspicam/picamera"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func (a *api) capture(w http.ResponseWriter, r *http.Request) {
	var imageSetKey key = "imageSetConfig"
	imageSetConfig := r.Context().Value(imageSetKey).(*picamera.ImageSetConfig)
	imageSet := imageSetConfig.Generate()
	if !a.captureInProgress {
		go func() {
			a.captureInProgress = !a.captureInProgress
			if err := imageSet.Capture(); err != nil {
				a.captureErr = err
			}
			a.captureInProgress = !a.captureInProgress
		}()
	} else {
		render.Render(w, r, &errResponse{
			HTTPStatusCode: 503,
			StatusText:     "Service Unavailable.",
			ErrorText:      "Camera sensor is busy",
		})
		return
	}

	w.Write([]byte("Capturing..."))
}

/*********
Middleware
*********/

type key string

func (a *api) loadImageSetConfig(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		var imageSetConfig *picamera.ImageSetConfig
		if imageSetConfigID := chi.URLParam(r, "imageSetConfigID"); imageSetConfigID != "" {
			imageSetConfig, err = getImageSetConfigFile(imageSetConfigID)
		} else {
			render.Render(w, r, errNotFound())
			return
		}
		if err != nil {
			render.Render(w, r, errNotFound())
			return
		}

		var imageSetKey key = "imageSetConfig"
		ctx := context.WithValue(r.Context(), imageSetKey, imageSetConfig)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

/***************
Helper functions
***************/

func getImageSetConfigFile(id string) (sc *picamera.ImageSetConfig, err error) {
	p, err := picamera.ImageSetConfigsPath()
	if err != nil {
		return
	}
	fpath := filepath.Join(p, id)
	file, err := ioutil.ReadFile(fpath)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(file), &sc)

	return
}
