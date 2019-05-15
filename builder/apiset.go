package builder

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/xerrors"
)

type (
	//APISet for multiple APIs
	APISet struct {
		APIs []API `json:"apis"`
	}
	// API in a set for multiple Commands
	API struct {
		Name     string
		BaseURL  string
		Header   map[string][]string
		Query    map[string][]string
		Commands []Command
		Set      *APISet
	}
	// Command on an API
	Command struct {
		Name   string
		Path   string
		Method string
		Header map[string][]string
		Query  map[string][]string
		API    *API
	}
)

func joinUrlFragments(base, path string) string {
	if base == "" {
		return path
	}
	if path == "" {
		return base
	}
	b := base[len(base)-1] == '/'
	p := path[0] == '/'
	if !b && !p {
		return base + "/" + path
	}
	if b && p {
		return base + path[1:]
	}
	return base + path
}

func validateURL(uri string) error {
	u, err := url.Parse(uri)
	if err != nil {
		return err
	}
	if u.Scheme != "https" && u.Scheme != "http" {
		return xerrors.Errorf("unsupported scheme: '%s'", uri)
	}
	if u.RawQuery != "" {
		return xerrors.Errorf("base URL cannot have query: '%s'", uri)
	}
	return nil
}

func (set *APISet) Prepare() {
	for _, api := range set.APIs {
		api.prepare()
	}
}

type (
	//NotFound api or command
	NotFound string
)

func (e NotFound) Error() string {
	return fmt.Sprintf("search for '%s' returned no results", string(e))
}

func (set *APISet) FindAPI(name string) ([]*API, error) {

	var list []*API
	var err error = NotFound(name)
	lowerName := strings.ToLower(name)
	for i := range set.APIs {
		//todo: can we do this WITHOUT the extra allocation?
		if strings.Contains(strings.ToLower(set.APIs[i].Name), lowerName) {
			list = append(list, &set.APIs[i])
			err = nil
		}
	}
	return list, err
}

func (api *API) FindCommand(name string) ([]*Command, error) {
	var list []*Command
	var err error = NotFound(name)
	lowerName := strings.ToLower(name)
	for i := range api.Commands {
		//todo: can we do this WITHOUT the extra allocation?
		if strings.Contains(strings.ToLower(api.Commands[i].Name), lowerName) {
			list = append(list, &api.Commands[i])
			err = nil
		}
	}
	return list, err
}

func (api *API) prepare() {
	for i := range api.Commands {
		api.Commands[i].API = api
	}
}
