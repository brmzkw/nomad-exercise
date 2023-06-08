package main

import (
	"fmt"
	"os"

	"github.com/brmzkw/nomad-exercise/internal/nomad"
	"github.com/brmzkw/nomad-exercise/internal/webservice"
)

func main() {
	nomad, err := nomad.NewNomad()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ws := webservice.NewWebService(nomad)
	ws.Run()
}
