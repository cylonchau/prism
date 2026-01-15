package main

import (
	"os"

	"github.com/cylonchau/prism/pkg/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
