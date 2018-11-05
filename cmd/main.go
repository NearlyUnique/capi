package main

import (
	"fmt"
	"log"
	"net/http"
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

	api, err := profile.SelectAPI(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
		return
	}
	cmd, err := profile.SelectCommand(api, os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
		return
	}
	fmt.Printf("%s :: %s\n\t%s\n", api.Name, cmd.Name, cmd.Path)
	fmt.Printf("\targs:\n\t%s\n", strings.Join(cmd.ListParams(), "\n"))
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
