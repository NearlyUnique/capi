package capi

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
)

func indexOrEmpty(args []string, i int) string {
	if i < 0 || i >= len(args) {
		return ""
	}
	return args[i]
}

//Prepare a command for given cli args
func (apiSet APISet) Prepare(args []string) (*Command, error) {
	api, err := apiSet.SelectAPI(indexOrEmpty(args, 0))
	if err != nil {
		return nil, err
	}
	cmd, err := api.SelectCommand(indexOrEmpty(args, 1))
	if err != nil {
		return nil, err
	}
	local := LocalEnv{
		EnvPrefix: apiSet.EnvPrefix,
		Environ:   os.Environ,
	}

	cmd.args, err = cmd.ParseArgs(args, local.Lookup)
	if err != nil {
		return nil, err
	}
	//err = fs.Parse(args[2:])
	//m := make(map[string]string)
	//cmd.fs.VisitAll(func(flag *pflag.Flag) {
	//	m[flag.Name] = flag.Value.String()
	//})

	//if err != nil {
	//	return nil, err
	//}
	return cmd, nil
}

//CreateRequest for sending, printing, debugging, etc.
func CreateRequest(cmd *Command) (*http.Request, error) {
	replace := func(src string) string {
		return rxMustacheParams.ReplaceAllStringFunc(src, func(s string) string {
			// s includes the '{' and '}'
			if v, ok := cmd.args[s[1:len(s)-1]]; ok {
				return v
			}
			return s
		})
	}

	tempPath := replace(cmd.Path)

	urlPath, err := makeURL(cmd.api.BaseURL, tempPath)
	if err != nil {
		return nil, err
	}
	var r io.Reader
	if cmd.Data != nil {
		applied := rxMustacheParams.ReplaceAllFunc(cmd.Data, func(b []byte) []byte {
			if v, ok := cmd.args[string(b[1:len(b)-1])]; ok {
				return []byte(v)
			}
			return b
		})
		r = bytes.NewReader(applied)
	}

	req, err := http.NewRequest(cmd.Method, urlPath, r)
	if err != nil {
		return nil, err
	}

	for k, v := range cmd.api.DefaultHeader {
		v = replace(v)
		req.Header.Set(k, v)
	}

	for k, v := range cmd.Header {
		v = replace(v)
		req.Header.Set(k, v)
	}

	return req, nil
}

func makeURL(baseURL, pathSegment string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	ps, err := url.Parse(pathSegment)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(u.Path, ps.Path)
	u.RawQuery = ps.Query().Encode()
	return u.String(), nil
}
