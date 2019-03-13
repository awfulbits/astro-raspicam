package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/awfulbits/astro-raspicam/config"
	"github.com/awfulbits/astro-raspicam/picamera"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type key string

const imageSetKey key = "imageSetConfig"

// imageSetConfigR struct is used for both request and response
type imageSetConfigR struct {
	ImageSetConfig picamera.ImageSetConfig `json:"imageSetConfig"`
}

func (sr *imageSetConfigR) Bind(r *http.Request) error {
	return nil
}

func (sr imageSetConfigR) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *server) saveImageSetConfig(w http.ResponseWriter, r *http.Request) {
	var err error
	var imageSetConfigR imageSetConfigR
	if err = render.Bind(r, &imageSetConfigR); err != nil {
		render.Render(w, r, errInvalidRequest(err))
		return
	}

	imageSetConfigR, err = saveImageSetConfigFile(imageSetConfigR.ImageSetConfig)
	if err != nil {
		render.Render(w, r, errInternalServer(err))
		return
	}

	render.Render(w, r, imageSetConfigR)
}

func (s *server) capture(w http.ResponseWriter, r *http.Request) {
	imageSetConfig := r.Context().Value(imageSetKey).(*picamera.ImageSetConfig)
	imageSet := imageSetConfig.Generate()
	if !s.captureInProgress {
		go func() {
			s.captureInProgress = !s.captureInProgress
			if err := imageSet.Capture(); err != nil {
				s.captureErr = err
			}
			s.captureInProgress = !s.captureInProgress
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

func (s *server) getImageSetConfig(w http.ResponseWriter, r *http.Request) {
	imageSetConfig := r.Context().Value(imageSetKey).(picamera.ImageSetConfig)
	render.Render(w, r, &imageSetConfigR{ImageSetConfig: imageSetConfig})
}

func (s *server) getImageSetConfigs(w http.ResponseWriter, r *http.Request) {
	scs, err := getImageSetConfigs()
	if err != nil {
		render.Render(w, r, errInternalServer(err))
		return
	}

	err = render.RenderList(w, r, scs)
	if err != nil {
		render.Render(w, r, errInternalServer(err))
		return
	}
}

/*********
Middleware
*********/

type imageSetConfigID struct {
	ImageSetConfigID string `json:"imageSetConfigID,omitempty"`
}

func (sr *imageSetConfigID) Bind(r *http.Request) error {
	return nil
}

func loadImageSetConfig(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		var imageSetConfigID imageSetConfigID
		if id := chi.URLParam(r, "imageSetConfigID"); id != "" {
			imageSetConfigID.ImageSetConfigID = id
		} else {
			if err = render.Bind(r, &imageSetConfigID); err != nil {
				render.Render(w, r, errInvalidRequest(err))
				return
			}
		}

		imageSetConfig, err := getImageSetConfigFile(imageSetConfigID.ImageSetConfigID)
		ctx := context.WithValue(r.Context(), imageSetKey, imageSetConfig)
		if err != nil {
			render.Render(w, r, errInternalServer(err))
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

/***************
Helper functions
***************/

func saveImageSetConfigFile(sc picamera.ImageSetConfig) (scr imageSetConfigR, err error) {
	if sc.ID == "" {
		seededRand := rand.New(
			rand.NewSource(time.Now().UnixNano()))
		sc.ID = strconv.Itoa(seededRand.Int())
	}

	scJSON, err := json.Marshal(&sc)
	fpath := filepath.Join(config.C.ImageSetConfigsPath, sc.ID)
	err = ioutil.WriteFile(fpath, scJSON, 0644)

	scr.ImageSetConfig = sc

	return
}

func getImageSetConfigFile(id string) (sc picamera.ImageSetConfig, err error) {
	fpath := filepath.Join(config.C.ImageSetConfigsPath, id)
	file, err := ioutil.ReadFile(fpath)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(file), &sc)

	return
}

func getImageSetConfigs() (scs []render.Renderer, e error) {
	root := config.C.ImageSetConfigsPath
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			sc, err := getImageSetConfigFile(info.Name())
			if err != nil {
				e = err
			}
			if sc.ID == info.Name() {
				scs = append(scs, imageSetConfigR{ImageSetConfig: sc})
			}
		}

		return nil
	})
	if err != nil {
		e = err
	}

	return
}
