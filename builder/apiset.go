package builder

import (
	"net/http"
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
		Header map[string][]string
		Query  map[string][]string
		API    *API
		Method string
	}
)

func (api *API) Add(cmd Command) error {
	if cmd.Name == "" {
		return xerrors.New("missing name")
	}
	cmd.API = api
	if cmd.Method == "" {
		cmd.Method = http.MethodGet
	}
	api.Commands = append(api.Commands, cmd)
	return nil
}

func (cmd *Command) CreateRequest() (*http.Request, error) {
	uri, _ := url.Parse(cmd.API.BaseURL + cmd.Path)
	req := &http.Request{
		Method: cmd.Method,
		URL:    uri,
	}
	return req, nil
}

func (set *APISet) Add(api API) error {
	if api.Name == "" {
		return xerrors.Errorf("missing name")
	}
	if api.BaseURL == "" {
		return xerrors.Errorf("missing baseUrl")
	}
	api.Set = set
	set.APIs = append(set.APIs, api)
	return nil
}
