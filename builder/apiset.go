package builder

import (
	"regexp"
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
		Header   map[string]StringOrList //[]string
		Query    map[string]StringOrList //[]string
		Commands []Command
		Set      *APISet
	}
	// Command on an API
	Command struct {
		Name   string
		Path   string
		Method string
		Header map[string]StringOrList //[]string
		Query  map[string]StringOrList //[]string
		Body   *CommandBody
		API    *API
	}
)

func (set *APISet) FindAPI(name string) ([]*API, error) {
	var list []*API
	var err error = NotFound(name)

	if len(set.APIs) == 1 {
		// we are here because a set was loaded
		// therefore if there is only one API
		// use it
		list = []*API{&set.APIs[0]}
		err = nil
	} else {
		lowerName := strings.ToLower(name)
		for i := range set.APIs {
			if strings.Contains(strings.ToLower(set.APIs[i].Name), lowerName) ||
				strings.ToLower(SafeName(set.APIs[i].Name)) == lowerName {
				api := &set.APIs[i]
				api.Set = set
				list = append(list, api)
				err = nil
			}
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
		if strings.Contains(strings.ToLower(api.Commands[i].Name), lowerName) ||
			strings.ToLower(SafeName(api.Commands[i].Name)) == lowerName {
			cmd := &api.Commands[i]
			cmd.API = api
			list = append(list, cmd)
			err = nil
		}
	}
	return list, err
}

func (api *API) FindCommandExact(name string) (*Command, error) {
	list, err := api.FindCommand(name)
	if err != nil {
		return nil, err
	}
	if len(list) == 1 {
		return list[0], nil
	}
	for _, c := range list {
		if strings.EqualFold(c.Name, name) {
			return c, nil
		}
	}
	return nil, NotFound(name)
}
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

var rxSafeName = regexp.MustCompile(`[\s()\[\]#\\/!"£$%^&*,<>;:@']`)

func SafeName(s string) string {
	return rxSafeName.ReplaceAllString(s, "_")
}
