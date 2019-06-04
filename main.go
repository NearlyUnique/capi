package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/NearlyUnique/capi/run"
)

func single(s string) []string {
	return []string{s}
}

func main() {
	homeDir, _ := os.UserHomeDir() // if can't read home folder , just carry on, might still work
	loader := run.NewConfigLoader(homeDir, ioutil.ReadFile)

	err := run.Main(loader, os.Args, os.Environ())
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error loading configuration: %v", err)
		os.Exit(1)
	}
}
