package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/NearlyUnique/capi/postman"
	"github.com/NearlyUnique/capi/run"
)

func main() {
	homeDir, _ := os.UserHomeDir() // if can't read home folder , just carry on, might still work

	loader := run.NewConfigLoader(homeDir, ioutil.ReadFile)
	loader.RegisterFileExtension(".json", run.JSONFormatReader)
	loader.RegisterFileExtension(".postman_collection.json", postman.FormatReader)

	if options := run.AutoComplete(loader, os.Args, os.Environ()); len(options) > 0 {
		for _, opt := range options {
			fmt.Println(opt)
		}
		return
	}

	err := run.Main(loader, os.Args, os.Environ())
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
}
