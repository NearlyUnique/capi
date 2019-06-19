package run

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/NearlyUnique/capi/autocomplete"
	"github.com/NearlyUnique/capi/builder"
	"github.com/NearlyUnique/capi/capicomplete"
)

// Main entry point
func Main(loader ConfigLoader, args, env []string) error {
	// find the config
	firstArg := indexOrEmpty(args, 1)
	set, err := loader.Load(firstArg)
	if err != nil {
		return err
	}

	req, err := set.CreateRequest(
		firstArg,
		indexOrEmpty(args, 2),
		builder.NewFlagSource(args, nil),
		builder.NewEnvVarsSource(env),
	)

	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	buf, _ := httputil.DumpResponse(resp, true)
	fmt.Print(string(buf))
	return nil
}

// AutoComplete entry point
func AutoComplete(loader ConfigLoader, args, env []string) []string {
	enableLogging(env)
	ac := autocomplete.Parse(args, env)

	if ac == nil {
		// ok, no complete work to perform
		return nil
	}
	// args are from complete, replace with the actual args that would run
	args = strings.Split(ac.Line, " ")
	// find the config
	firstArg := indexOrEmpty(args, 1)
	set, err := loader.Load(firstArg)
	if err != nil {
		return loader.List(firstArg)
	}

	return capicomplete.GenerateResponse(ac, set)
}

// enableLogging is performed if CAPI_DEBUG has any value
func enableLogging(env []string) {
	for _, v := range env {
		if strings.HasPrefix(v, "CAPI_DEBUG") {
			f, err := os.Create("capi.log")
			if err != nil {
				autocomplete.LogHook = log.Printf
			} else {
				autocomplete.LogHook = func(format string, args ...interface{}) {
					_, _ = fmt.Fprintf(f, format, args...)
				}
			}
			return
		}
	}
}

type ConfigLoader interface {
	List(search string) []string
	Load(filename string) (*builder.APISet, error)
}
