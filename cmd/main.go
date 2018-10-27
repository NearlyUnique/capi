package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/NearlyUnique/capi/capicomplete"

	"github.com/NearlyUnique/capi"
	"github.com/NearlyUnique/capi/autocomplete"
)

func main() {
	printVersion(os.Args)
	profilePath := "./profile.json"
	if _, err := os.Stat(profilePath); os.IsNotExist(err) {
		profilePath = os.Getenv("CAPI_PROFILE")
	}
	if profilePath == "" {
		fmt.Fprintln(os.Stderr, "./profile.json does not exist and ENV var CAPI_PROFILE not set")
		os.Exit(1)
	}
	f, err := os.Open(profilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't open profile %s\n", profilePath)
		os.Exit(1)
	}
	profile, err := capi.LoadAPI(f)

	ac := autocomplete.Prepare(os.Args, os.Environ())
	if ac != nil {
		comp := capicomplete.GenerateResponse(ac, profile)
		fmt.Print(strings.Join(comp, "\n"))
		return
	}

	//api, err := whichApi(os.Args[1:])
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "FATAL: %v\n", err)
	//	os.Exit(1)
	//}
	//cmd, err := api.whichCmd(os.Args[2:])
	//
	//paramNames := extractParams(api, cmd)
	//values := make(map[string]*string)
	//
	//fs := flag.NewFlagSet(api.Alias, flag.ContinueOnError)
	//for _, p := range paramNames {
	//	values[p] = fs.String(p, "no description", "optional value")
	//}
	//
	//// 0        1        2        3
	//// this-app api-Name cmd-Name firstarg
	//fs.Parse(os.Args[3:])
	//fmt.Println(fs.Parsed())
	//
	//h := applyParams(api, cmd, values)
	//
	//err = httpRequest(api, cmd, h)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "FATAL:%v\n", err)
	//	os.Exit(1)
	//}
}

var rxMustacheParams = regexp.MustCompile(`{(?P<Name>[a-zA-Z0-9-_]+)}`)

func applyParams(api *capi.API, cmd *capi.Command, values map[string]*string) http.Header {

	//// from url
	//cmd.Path = rxMustacheParams.ReplaceAllStringFunc(cmd.Path, func(s string) string {
	//	k := s[1 : len(s)-1]
	//	v, ok := values[k]
	//	if ok && *v != "" {
	//		return *v
	//	}
	//	return s
	//})
	//// from host
	//host := api.BaseURL[config("env", "dev")]
	//_ = rxMustacheParams.ReplaceAllStringFunc(host, func(s string) string {
	//	return s + "###"
	//})

	// from headers
	h := make(http.Header)
	//for _, p := range cmd.Header {
	//	if v, ok := values[p]; ok {
	//		h.Set(p, *v)
	//	}
	//}
	return h
}
func extractParams(api *capi.API, cmd *capi.Command) []string {
	var params []string
	// from url
	//match := rxMustacheParams.FindAllStringSubmatch(cmd.Path, -1)
	//for _, m := range match {
	//	params = append(params, m[1])
	//}
	//// from host
	//match = rxMustacheParams.FindAllStringSubmatch(api.BaseURL[config("env", "dev")], -1)
	//for _, m := range match {
	//	params = append(params, m[1])
	//}
	//
	//// from headers
	//params = append(params, cmd.Header...)

	return params
}
func httpRequest(api *capi.API, cmd *capi.Command, header http.Header) error {
	//c := http.Client{}
	//req, err := http.NewRequest(cmd.Method, api.BaseURL(cmd), nil)
	//if err != nil {
	//	return err
	//}
	//req.Header = header
	//
	//resp, err := c.Do(req)
	//if err != nil {
	//	return err
	//}
	//
	//defer resp.Body.Close()
	//
	//for k, v := range resp.Header {
	//	for _, h := range v {
	//		fmt.Printf("%v: %v\n", k, h)
	//	}
	//}
	//
	//_, err = io.Copy(os.Stdout, resp.Body)
	//if err != nil {
	//	return err
	//}

	return nil
}

//func (a *capi.API) url(cmd *capi.Command) string {
//	env := config("env", "dev")
//	baseURL, ok := a.BaseURL[env]
//	if !ok {
//		fmt.Fprintf(os.Stderr, "configured environment %s not found", env)
//		os.Exit(1)
//	}
//	u, _ := url.Parse(baseURL)
//	u.Path = path.Join(u.Path, cmd.Path)
//	return u.String()
//}
//func config(key, def string) string {
//	if v, ok := privateEnv[key]; ok {
//		return v
//	}
//	return def
//}

//func (a *API) whichCmd(args []string) (*Command, error) {
//	if len(args) < 1 {
//		return nil, errors.Errorf("no cmd selected")
//	}
//	c, ok := a.Commands[args[0]]
//	if !ok {
//		return nil, errors.Errorf("no api named %s registered", args)
//	}
//	return c, nil
//
//}
//func whichApi(args []string) (*API, error) {
//	if len(args) < 1 {
//		return nil, errors.Errorf("no api selected")
//	}
//	a, ok := apis[args[0]]
//	if !ok {
//		return nil, errors.Errorf("no api named %s registered", args)
//	}
//	return &a, nil
//}

//func init() {
//	for _, pair := range os.Environ() {
//		s := strings.Split(pair, "=")
//		key := strings.ToLower(s[0])
//		if strings.HasPrefix(key, "crv_") {
//			privateEnv[key[:4]] = s[1]
//		}
//	}
//}

var version = "0.0"

func printVersion(args []string) {
	if len(args) >= 2 && args[1] == "version" {
		fmt.Println(version)
		os.Exit(0)
	}
}
