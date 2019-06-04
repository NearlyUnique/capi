package builder

import (
	"github.com/NearlyUnique/pflag"
)

func NewFlagSource(args []string, errLog SourceErrorFn) SourceFn {
	flagSet := pflag.NewFlagSet("any", pflag.ContinueOnError)
	flagSet.ParseErrorsWhitelist.UnknownFlags = true

	err := flagSet.Parse(args)
	if err != nil {
		errLog(err)
	}

	return func(k string) string {
		uk, ok := flagSet.UnknownFlags[k]
		if !ok {
			return ""
		}
		return uk[0].Value
	}
}
