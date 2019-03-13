package main

import (
	"bytes"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func (s *server) getRoot(w http.ResponseWriter, r *http.Request) {
	if err := s.tmpls.ExecuteTemplate(w, "layout", nil); err != nil {
		render.Render(w, r, errInternalServer(err))
	}
}

func (s *server) getModule(w http.ResponseWriter, r *http.Request) {
	var tpl bytes.Buffer

	if err := s.tmpls.ExecuteTemplate(&tpl, chi.URLParam(r, "tmpl"), nil); err != nil {
		render.Render(w, r, errInternalServer(err))
	}

	w.Write(tpl.Bytes())
}
