package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/NearlyUnique/capi"
	"github.com/NearlyUnique/capi/autocomplete"
	"github.com/NearlyUnique/capi/capicomplete"
)

func main() {
	printVersion(os.Args)
	if err := run(); err != nil {
		logError(err)
		os.Exit(1)
	}
}

func run() error {
	apiSet, err := capi.Load(capi.DefaultAPISet)
	if err != nil {
		return err
	}
	if autoComplete(apiSet) {
		return nil
	}
	if len(os.Args) == 1 {
		return errors.New("Insufficient arguments")
	}
	cmd, err := apiSet.Prepare(os.Args[1:])
	if err != nil {
		return err
	}
	req, err := capi.CreateRequest(cmd)
	if err != nil {
		return err
	}
	// make json if --capi-mode = request or all
	//b, _ := httputil.DumpRequest(req, true)
	//fmt.Println(string(b))

	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	// make json if --capi-mode = response or all
	//b, _ = httputil.DumpResponse(resp, true)
	//fmt.Println(string(b))
	if resp.Body != nil {
		defer func() {
			logError(resp.Body.Close())
		}()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		fmt.Print(string(b))
	}
	return nil
}

func logError(err error) {
	if err == nil {
		return
	}
	// if you can't write to srd error, no one can read the log
	_, _ = fmt.Fprintln(os.Stderr, err.Error())
}

func autoComplete(apiSet *capi.APISet) bool {
	autocomplete.LogHook = logFn
	ac := autocomplete.Parse(os.Args, os.Environ())
	if ac != nil {
		logFn("%v", ac)
		comp := capicomplete.GenerateResponse(ac, apiSet)
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
		logError(err)
	}
}

//func httpRequest(api *capi.API, cmd *capi.Command, header http.Header) error {
//	c := http.Client{}
//	req, err := http.NewRequest(cmd.Method, "", nil) //api.BaseURL(cmd), nil)
//	if err != nil {
//		return err
//	}
//	req.Header = header
//
//	resp, err := c.Do(req)
//	if err != nil {
//		return err
//	}
//
//	defer resp.Body.Close()
//
//	for k, v := range resp.Header {
//		for _, h := range v {
//			fmt.Printf("%v: %v\n", k, h)
//		}
//	}
//
//	_, err = io.Copy(os.Stdout, resp.Body)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

var version = "0.0"

func printVersion(args []string) {
	if len(args) >= 2 && args[1] == "version" {
		fmt.Println(version)
		os.Exit(0)
	}
}
