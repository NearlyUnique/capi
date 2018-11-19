package capi

import (
	"net/http"
	"net/url"
	"os"
	"path"
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

	cmd.args, err = cmd.ParseArgs(args, local.Lookup)
	if err != nil {
		return nil, err
	}
	//err = fs.Parse(args[2:])
	//m := make(map[string]string)
	//cmd.fs.VisitAll(func(flag *pflag.Flag) {
	//	m[flag.Name] = flag.Value.String()
	//})

	if err != nil {
		return nil, err
	}
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

	urlPath := makeURL(cmd.api.BaseURL, tempPath)

	req, err := http.NewRequest(cmd.Method, urlPath, nil)
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

func makeURL(baseURL, pathSegment string) string {
	u, _ := url.Parse(baseURL)
	u.Path = path.Join(u.Path, pathSegment)
	return u.String()
}
