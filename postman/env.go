package postman

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"

	"github.com/NearlyUnique/capi/builder"
)

const (
	CollectionExtension  = ".postman_collection.json"
	EnvironmentExtension = ".postman_environment.json"
)

type (
	FileOpener func(string) (io.Reader, error)
	Config     struct {
		args, env []string
		opener    FileOpener
		errFn     builder.SourceErrorFn
	}
)

func NewSource(args, env []string, reader FileOpener, fn builder.SourceErrorFn) builder.SourceFn {
	if fn == nil {
		fn = func(error) {}
	}
	return Config{
		args:   args,
		env:    env,
		opener: reader,
		errFn:  fn,
	}.Source()
}

func (c Config) Source() builder.SourceFn {
	buf, err := c.loadEnv(resolveFilename(c.args, c.env))
	if err != nil {
		c.errFn(err)
		return nil
	}
	m, err := envFileToMap(buf)
	if err != nil {
		c.errFn(err)
		return nil
	}
	return builder.SourceFn(func(key string) string {
		return m[key]
	})
}

func (c Config) loadEnv(filename string) ([]byte, error) {
	if !strings.HasSuffix(filename, EnvironmentExtension) {
		filename += EnvironmentExtension
	}
	f, err := c.opener(filename)
	if err != nil {
		return nil, err
	}
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func envFileToMap(buf []byte) (map[string]string, error) {
	var postCfg struct {
		Values []struct {
			Key, Value string
		}
	}
	err := json.Unmarshal(buf, &postCfg)
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)
	for _, kv := range postCfg.Values {
		m[kv.Key] = kv.Value
	}
	return m, nil
}

func resolveFilename(args, env []string) string {
	const postmanEnv = "CAPI_ENV="
	for _, e := range env {
		if strings.HasPrefix(e, postmanEnv) {
			val := e[len(postmanEnv):]
			if len(val) > 0 {
				return val
			}
		}
	}
	if len(args) > 1 {
		return args[0]
	}
	return ""
}
