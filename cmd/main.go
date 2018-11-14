package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"
	"strings"

	"github.com/NearlyUnique/capi"
	"github.com/NearlyUnique/capi/autocomplete"
	"github.com/NearlyUnique/capi/capicomplete"
	"github.com/pkg/errors"
)

func main() {
	printVersion(os.Args)
	profile, err := loadProfile()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	if autoComplete(profile) {
		return
	}

	cmd, err := capi.Prepare(*profile, os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	req, err := capi.CreateRequest(cmd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	b, _ := httputil.DumpRequest(req, true)
	fmt.Println(string(b))

	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	b, _ = httputil.DumpResponse(resp, true)
	fmt.Println(string(b))
}

func loadProfile() (*capi.Profile, error) {
	profilePath := "./profile.json"
	if _, err := os.Stat(profilePath); os.IsNotExist(err) {
		profilePath = os.Getenv("CAPI_PROFILE")
	}
	if profilePath == "" {
		return nil, errors.New("./profile.json does not exist and ENV var CAPI_PROFILE not set")
	}
	f, err := os.Open(profilePath)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("can't open profile %s\n", profilePath))
	}
	return capi.LoadAPI(f)
}

func autoComplete(profile *capi.Profile) bool {
	autocomplete.LogHook = logFn
	ac := autocomplete.Prepare(os.Args, os.Environ())
	if ac != nil {
		logFn("%v", ac)
		comp := capicomplete.GenerateResponse(ac, profile)
		fmt.Print(strings.Join(comp, "\n"))
		return true
	}
	return false
}

func logFn(format string, args ...interface{}) {
	logfile, err := os.OpenFile("complete.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
	if err == nil {
		defer logfile.Close()

		log.SetOutput(logfile)
		log.Printf(format, args...)
	} else {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

var rxMustacheParams = regexp.MustCompile(`{(?P<Name>[a-zA-Z0-9-_]+)}`)

func httpRequest(api *capi.API, cmd *capi.Command, header http.Header) error {
	c := http.Client{}
	req, err := http.NewRequest(cmd.Method, "", nil) //api.BaseURL(cmd), nil)
	if err != nil {
		return err
	}
	req.Header = header

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	for k, v := range resp.Header {
		for _, h := range v {
			fmt.Printf("%v: %v\n", k, h)
		}
	}

	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

var version = "0.0"

func printVersion(args []string) {
	if len(args) >= 2 && args[1] == "version" {
		fmt.Println(version)
		os.Exit(0)
	}
}
