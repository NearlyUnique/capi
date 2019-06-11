package builder

import (
	"io"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/xerrors"
)

//CreateRequest by searching the set for a single api/command combination
func (set *APISet) CreateRequest(api, command string, sources ...SourceFn) (*http.Request, error) {
	if set == nil {
		return nil, InvalidOperation("nil apiset")
	}
	apis, err := set.FindAPI(api)
	if err != nil {
		return nil, err
	}
	if len(apis) != 1 {
		return nil, NotFound(api)
	}
	cmds, err := apis[0].FindCommand(command)
	if err != nil {
		return nil, err
	}
	if len(cmds) != 1 {
		return nil, NotFound(command)
	}
	cmds[0].API = apis[0]
	return cmds[0].CreateRequest(sources...)
}

// CreateRequest
func (cmd Command) CreateRequest(sources ...SourceFn) (*http.Request, error) {
	if cmd.API == nil {
		return nil, xerrors.New("APISet.Prepare has not been called")
	}
	uri := joinUrlFragments(cmd.API.BaseURL, cmd.Path)
	applier := func(value string) string {
		return applyReplacement(value, sources)
	}
	uri = applier(uri)

	var payload io.Reader
	if cmd.Body != nil {
		buf := applier(string(cmd.Body.Data))
		payload = strings.NewReader(buf)
	}
	req, err := http.NewRequest(cmd.Method, uri, payload)
	if err != nil {
		return nil, err
	}

	applyToHeaders(cmd.API.Header, req, applier)
	applyToHeaders(cmd.Header, req, applier)

	q := req.URL.Query()
	for k, v := range cmd.API.Query {
		for _, val := range v {
			q.Add(k, applier(val))
		}
	}
	for k, v := range cmd.Query {
		q.Del(k)
		for _, val := range v {
			q.Add(k, applier(val))
		}
	}
	req.URL.RawQuery = q.Encode()

	return req, nil
}

func applyToHeaders(headers map[string]StringOrList, req *http.Request, apply func(v string) string) {
	for k, v := range headers {
		req.Header.Del(k)
		for _, item := range v {
			req.Header.Add(k, apply(item))
		}
	}
}

var rxMustacheParams = regexp.MustCompile(`{(?P<Name>[a-zA-Z0-9-_]+)}`)

func applyReplacement(source string, sources []SourceFn) string {
	return replaceAllStringSubmatchFunc(rxMustacheParams, source, func(args []string) string {
		for _, srcFn := range sources {
			if v := srcFn(args[1]); v != "" {
				return v
			}
		}
		return args[0]
	})
}

func replaceAllStringSubmatchFunc(re *regexp.Regexp, str string, repl func(args []string) string) string {
	result := ""
	lastIndex := 0

	for _, v := range re.FindAllSubmatchIndex([]byte(str), -1) {
		var groups []string
		for i := 0; i < len(v); i += 2 {
			groups = append(groups, str[v[i]:v[i+1]])
		}

		result += str[lastIndex:v[0]] + repl(groups)
		lastIndex = v[1]
	}

	return result + str[lastIndex:]
}
