package main

import (
	"os"

	"github.com/cylonchau/prism/cmd/prism/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
