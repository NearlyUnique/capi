package capi

import (
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/spf13/pflag"
)

//Prepare a command for given cli args
func Prepare(profile Profile, args []string) (*Command, error) {
	api, err := profile.SelectAPI(args)
	if err != nil {
		return nil, err
	}
	cmd, err := profile.SelectCommand(api, args)
	if err != nil {
		return nil, err
	}
	local := LocalEnv{
		EnvPrefix: profile.EnvPrefix,
		Environ:   os.Environ,
	}

	cmd.fs = cmd.CreateFlagSet(local.Lookup)
	err = cmd.fs.Parse(args[2:])

	if err != nil {
		return nil, err
	}
	return cmd, nil
}

//CreateRequest for sending, printing, debugging, etc.
func CreateRequest(cmd *Command) (*http.Request, error) {
	m := make(map[string]string)
	cmd.fs.VisitAll(func(flag *pflag.Flag) {
		m[flag.Name] = flag.Value.String()
	})

	tempPath := rxMustacheParams.ReplaceAllStringFunc(cmd.Path, func(s string) string {
		if v, ok := m[s[1:len(s)-1]]; ok {
			return v
		}
		return s
	})

	urlPath := makeURL(cmd.api.BaseURL, tempPath)

	req, err := http.NewRequest(cmd.Method, urlPath, nil)
	if err != nil {
		return nil, err
	}
	//req.Header
	return req, nil
}

func makeURL(baseURL, pathSegment string) string {
	u, _ := url.Parse(baseURL)
	u.Path = path.Join(u.Path, pathSegment)
	return u.String()
}
