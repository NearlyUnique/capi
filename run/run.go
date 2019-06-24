package run

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/NearlyUnique/capi/autocomplete"
	"github.com/NearlyUnique/capi/builder"
	"github.com/NearlyUnique/capi/capicomplete"
)

type LoadLister interface {
	List(search string) []string
	Load(filename string) (*builder.APISet, error)
}

// Main entry point
func Main(loader LoadLister, args []string, sources ...builder.SourceFn) error {
	// find the config
	firstArg := indexOrEmpty(args, 1)
	set, err := loader.Load(firstArg)
	if err != nil {
		return err
	}

	req, err := set.CreateRequest(
		firstArg,
		indexOrEmpty(args, 2),
		sources...,
	)

	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	buf, err := json.Marshal(resp)
	if err != nil {
		buf, _ = httputil.DumpResponse(resp, true)
	}
	fmt.Print(string(buf))

	return nil
}

// AutoComplete entry point
func AutoComplete(loader LoadLister, args, env []string) []string {
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
