package run

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/NearlyUnique/capi/autocomplete"
	"github.com/NearlyUnique/capi/builder"
	"github.com/NearlyUnique/capi/capicomplete"
)

//Main entry point
func Main(loader ConfigLoader, args, env []string) error {
	loader.RegisterFileExtension(".json", JSONFormatReader)
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
	buf, _ := httputil.DumpResponse(resp, true)
	fmt.Print(string(buf))
	return nil
}

func AutoComplete(loader ConfigLoader, args, env []string) bool {
	loader.RegisterFileExtension(".json", JSONFormatReader)
	// find the config
	firstArg := indexOrEmpty(args, 1)
	set, err := loader.Load(firstArg)
	if err != nil {
		// log?
		return false
	}

	ac := autocomplete.Parse(args, env)

	if ac == nil {
		// ok, no complete work to perform
		return false
	}
	options := capicomplete.GenerateResponse(ac, set)
	if len(options) > 0 {
		for _, opt := range options {
			fmt.Println(opt)
		}
	}

	return true
}
