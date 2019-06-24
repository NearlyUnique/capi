package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/NearlyUnique/capi/autocomplete"
	"github.com/NearlyUnique/capi/builder"
	"github.com/NearlyUnique/capi/capicomplete"
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
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func osOpenFile(filename string) (reader io.Reader, e error) {
	return os.Open(filename)
}

func errorLog(e error) {
	log.Print(e.Error())
}

func init() {
	if name := os.Getenv("CAPI_LOG"); len(name) > 0 {
		var w io.Writer
		var err error
		w, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			// asked for logging but we could not so print to screen
			w = os.Stdout
			log.Printf("FAILED to create log file '%s': %v", name, err)
		}
		log.SetOutput(w)
		autocomplete.EnableLogging(w)
		postman.EnableLogging(w)
		run.EnableLogging(w)
		builder.EnableLogging(w)
		capicomplete.EnableLogging(w)
	}

}
