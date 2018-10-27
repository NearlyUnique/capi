package capi

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

type (
	apiMap  map[string]API
	Profile struct {
		EnvPrefix string `json:"envPrefix"`
		APIs      []API  `json:"apis"`
	}
	API struct {
		Name     string            `json:"name"`
		Alias    string            `json:"alias"`
		BaseURL  map[string]string `json:"baseURL"` // per environment (key)
		Commands []Command         `json:"commands"`
	}
	Command struct {
		Name   string   `json:"name"`
		Method string   `json:"method"`
		Path   string   `json:"path"`
		Header []string `json:"header"`
	}
)

func (cmd Command) CurlString(baseURL string) string {
	return fmt.Sprintf("# %s\ncurl -X %s %s%s", cmd.Name, cmd.Method, baseURL, cmd.Path)
}

var (
	apis       apiMap
	privateEnv map[string]string
)

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
