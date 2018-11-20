package capi

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

type (
	Profile struct {
		EnvPrefix string `json:"envPrefix"`
		APIs      []API  `json:"apis"`
	}
	API struct {
		Name          string            `json:"name"`
		BaseURL       string            `json:"baseURL"`
		DefaultHeader map[string]string `json:"defaultHeader"`
		Commands      []Command         `json:"commands"`

		profile *Profile
	}
	Command struct {
		Name   string            `json:"name"`
		Method string            `json:"method"`
		Path   string            `json:"path"`
		Header map[string]string `json:"header"`
		Data   json.RawMessage   `json:"data"`

		api  *API
		args map[string]string
	}
)

const (
	argIndexAPI     = 1
	argIndexCommand = 2
)

func (cmd Command) CurlString(baseURL string) string {
	return fmt.Sprintf("# %s\ncurl -X %s %s%s", cmd.Name, cmd.Method, baseURL, cmd.Path)
}

// LoadAPI from a basic config json file
func LoadAPI(r io.ReadCloser) (*Profile, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var profile Profile
	err = json.Unmarshal(buf, &profile)
	if err != nil {
		return nil, err
	}
	return &profile, err
}

// SelectAPI selects based on the os.Args
func (p *Profile) SelectAPI(args []string) (*API, error) {
	if len(args) <= argIndexAPI {
		return nil, errors.New("not enough arguments")
	}
	apiName := args[argIndexAPI]
	for _, a := range p.APIs {
		if apiName == a.Name {
			a.profile = p
			return &a, nil
		}
	}
	return nil, errors.Errorf("no api named %s registered", apiName)
}

//SelectCommand from an API
func (p *Profile) SelectCommand(api *API, args []string) (*Command, error) {
	if api == nil {
		return nil, errors.New("nil api")
	}
	if len(args) <= argIndexCommand {
		return nil, errors.New("not enough arguments")
	}
	cmdName := args[argIndexCommand]
	for _, c := range api.Commands {
		if cmdName == c.Name {
			api.Add(&c)
			return &c, nil
		}
	}
	return nil, errors.Errorf("no command named %s registered", cmdName)
}

//Add a cmd to the api
func (api *API) Add(cmd *Command) {
	if cmd.api == api {
		return
	}
	cmd.api = api
	for _, c := range api.Commands {
		if c.Method == cmd.Method && c.Path == cmd.Path {
			return
		}
	}

	api.Commands = append(api.Commands, *cmd)
}

var rxMustacheParams = regexp.MustCompile(`{(?P<Name>[a-zA-Z0-9-_]+)}`)

func (cmd Command) ListParams() []string {
	var params []string
	set := map[string]string{}
	add := func(v string) {
		if _, ok := set[v]; ok {
			return
		}
		set[v] = v
		params = append(params, v)

	}
	addMap := func(header map[string]string) {
		for k, v := range header {
			//values
			match := rxMustacheParams.FindAllStringSubmatch(v, -1)
			for _, m := range match {
				add(m[1])
			}
			// keys, only if no values inside
			if len(match) == 0 {
				add(k)
			}
		}
	}
	//from url
	match := rxMustacheParams.FindAllStringSubmatch(cmd.Path, -1)
	for _, m := range match {
		add(m[1])
	}
	// header
	addMap(cmd.Header)
	// from default headers
	if cmd.api != nil {
		addMap(cmd.api.DefaultHeader)
	}
	return params
}

// TODO: should there be any defaults?
// ParseArgs from the command line adding any defaults from env vars and profile
func (cmd *Command) ParseArgs(args []string, eval func(string) string) (map[string]string, error) {
	fs := pflag.NewFlagSet(cmd.Name, pflag.ContinueOnError)
	if eval == nil {
		eval = func(string) string { return "" }
	}
	for _, p := range cmd.ListParams() {
		fs.String(p, eval(p), "")
	}

	err := fs.Parse(args[2:])
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)
	fs.VisitAll(func(flag *pflag.Flag) {
		m[flag.Name] = flag.Value.String()
	})
	return m, nil
}

type LocalEnv struct {
	EnvPrefix string
	Environ   func() []string
	envCache  map[string]string
}

//Lookup a value from the environment
func (e *LocalEnv) Lookup(key string) string {
	if len(e.envCache) == 0 {
		e.envCache = make(map[string]string)
		for _, v := range e.Environ() {
			i := strings.Index(v, "=")
			if i == -1 {
				continue
			}
			e.envCache[strings.ToUpper(v[:i])] = v[i+1:]
		}
	}
	kup := strings.ToUpper(key)

	if v, ok := e.envCache[e.EnvPrefix+kup]; ok {
		return v
	}
	if v, ok := e.envCache[kup]; ok {
		return v
	}

	return ""
}
