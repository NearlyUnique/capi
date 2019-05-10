package builder

import (
	"net/http"
	"regexp"
)

// CreateRequest
func CreateRequest(cmd Command, sources ...SourceFn) (*http.Request, error) {
	uri := joinUrlFragments(cmd.API.BaseURL, cmd.Path)
	uri = applyReplacement(uri, sources)

	req, err := http.NewRequest(cmd.Method, uri, nil)
	if err != nil {
		return nil, err
	}

	applier := func(value string) string {
		return applyReplacement(value, sources)
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

func applyToHeaders(headers map[string][]string, req *http.Request, apply func(v string) string) {
	for k, v := range headers {
		//req.Header.Del(k)
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
