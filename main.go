package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/NearlyUnique/capi/builder"
	"github.com/NearlyUnique/capi/postman"
	"github.com/NearlyUnique/capi/run"
)

func main() {
	homeDir, _ := os.UserHomeDir() // if can't read home folder , just carry on, might still work

	loader := run.NewConfigLoader(homeDir, ioutil.ReadFile)
	loader.RegisterFileExtension(run.JSONFileExtension, run.JSONFormatReader)
	loader.RegisterFileExtension(postman.CollectionExtension, postman.FormatReader)

	env := os.Environ()
	if options := run.AutoComplete(loader, os.Args, env); len(options) > 0 {
		for _, opt := range options {
			fmt.Println(opt)
		}
		return
	}
	err := run.Main(loader, os.Args,
		builder.NewFlagSource(os.Args, errorLog),
		builder.NewEnvVarsSource(env),
		postman.NewSource(os.Args, env, osOpenFile, errorLog),
	)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
}

func osOpenFile(filename string) (reader io.Reader, e error) {
	return os.Open(filename)
}

func errorLog(e error) {

}
