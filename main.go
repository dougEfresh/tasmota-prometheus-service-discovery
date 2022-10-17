package main

import (
	"fmt"
	"os"

	"github.com/dougEfresh/tasmota-prometheus-service-discovery/cmd"
)

var version = ""

// nolint:revive
func main() {
	if err := cmd.Execute(version); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}
