package postman

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/NearlyUnique/capi/builder"
)

type (
	// Value represents the complex types from postman data, not all values are completed
	// in all cases
	Value struct {
		Key   string
		Value string
		Name  string
		Type  string
	}
	Info struct {
		PostmanID string `json:"_postman_id"`
		Name      string
		Schema    string
	}
	Auth struct {
		Type   string
		Bearer []Value
	}
	Body struct {
		Mode string
		Raw  string
	}
	URL struct {
		Raw      string
		Protocol string
		Host     []string
		Path     []string
		Query    []Value
	}
	Request struct {
		Description string
		Auth        Auth
		Method      string
		Header      []Value
		Body        Body
		URL         URL
	}
	Item struct {
		Name    string
		Request Request
		//Event []struct {
		//	Listen string
		//	Script struct {
		//		ID   string
		//		Exec []string
		//		Type string
		//	}
		//}
	}
	//ItemConfig from a postman file
	ItemConfig struct {
		Info Info
		Item []Item
	}
)

// FormatReader to read postman configuration and convert to capi config
func FormatReader(content []byte) (*builder.APISet, error) {
	var pm ItemConfig
	err := json.Unmarshal(content, &pm)
	if err != nil {
		return nil, err
	}

	return pm.ToAPISet()
}

// ToAPISet to builder.APISet
func (c ItemConfig) ToAPISet() (*builder.APISet, error) {
	set := builder.APISet{
		APIs: []builder.API{
			{Name: c.Info.Name},
		},
	}
	for _, i := range c.Item {
		cmd := builder.Command{
			Name:   strings.Replace(i.Name, " ", "_", -1),
			Path:   convertToCapiParameters(strings.Split(i.Request.URL.Raw, "?")[0]),
			Method: i.Request.Method,
			Header: make(map[string]builder.StringOrList),
			Query:  make(map[string]builder.StringOrList),
		}
		postmanListToMap(i.Request.Header, cmd.Header)
		postmanListToMap(i.Request.URL.Query, cmd.Query)

		if len(i.Request.Body.Raw) > 0 {
			s := i.Request.Body.Raw
			cmd.Body = &builder.CommandBody{Data: []byte(convertToCapiParameters(s))}
		}
		set.APIs[0].Commands = append(set.APIs[0].Commands, cmd)
	}

	return &set, nil
}

func postmanListToMap(item []Value, cmdMap map[string]builder.StringOrList) {
	for _, h := range item {
		val := convertToCapiParameters(h.Value)
		if _, ok := cmdMap[h.Key]; ok {
			cmdMap[h.Key] = append(cmdMap[h.Key], val)
		} else {
			cmdMap[h.Key] = []string{val}
		}
	}
}

var rxMustacheParams = regexp.MustCompile(`{{(?P<Name>[a-zA-Z0-9-_]+)}}`)

func convertToCapiParameters(postman string) string {
	return rxMustacheParams.ReplaceAllString(postman, "{$1}")
}
