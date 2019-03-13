package main

import (
	"html/template"
	"log"
	"net/http"
)

type server struct {
	tmpls             *template.Template
	captureInProgress bool
	captureErr        error
}

func main() {
	s := &server{}
	// Preload templates because theres no reason not to
	s.tmpls = template.Must(template.ParseGlob("templates/*"))

	log.Fatal(http.ListenAndServe(":3333", s.createRouter()))
}
