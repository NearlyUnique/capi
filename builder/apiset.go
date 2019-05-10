package builder

import (
	"net/url"

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
	for i, api := range set.APIs {
		set.APIs[i].Set = set
		api.prepare()
	}
}

func (api *API) prepare() {
	for i := range api.Commands {
		api.Commands[i].API = api
	}
}
