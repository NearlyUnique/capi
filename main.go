package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/NearlyUnique/capi/builder"
)

func single(s string) []string {
	return []string{s}
}

func main() {
	homeDir, _ := os.UserHomeDir() // if can't read home folder , just carry on, might still work
	loader := NewConfigLoader(homeDir, ioutil.ReadFile)
	loader.RegisterFileExtension(".json", JSONFormatReader)
	// find the config
	firstArg := indexOrEmpty(os.Args, 1)
	set, err := loader.Load(firstArg)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error loading configuration: %v", err)
		os.Exit(1)
		return
	}

	req, err := set.CreateRequest(
		indexOrEmpty(os.Args, 1),
		indexOrEmpty(os.Args, 2),
		builder.NewCaseInsensitiveSource(os.Environ()))

	if err != nil {
		os.Exit(1)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	buf, _ := httputil.DumpResponse(resp, true)
	fmt.Print(string(buf))
}
