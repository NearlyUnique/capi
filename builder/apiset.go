package builder

import (
	"strings"
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
		Body   *CommandBody
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

func (set *APISet) Prepare() {
	for _, api := range set.APIs {
		api.prepare()
	}
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
