package main

import (
	"log"

	"github.com/awfulbits/astro-raspicam/api"
)

func main() {
	log.Fatal(api.Start())
}
