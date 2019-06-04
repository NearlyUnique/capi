package builder

import (
	"strings"
)

type SourceFn func(string) string
type SourceErrorFn func(error)

func NewEnvVarsSource(list []string) SourceFn {
	m := make(map[string]string)
	for _, e := range list {
		parts := strings.Split(e, "=")
		m[parts[0]] = e[len(parts[0])+1:]
	}
	return SourceFn(func(key string) string {
		// literal match
		if s, ok := m[key]; ok {
			return s
		}
		// case insensitive search
		for k, v := range m {
			if strings.EqualFold(k, key) {
				return v
			}
		}
		return ""
	})
}
