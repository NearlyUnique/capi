package autocomplete

import (
	"strconv"
	"strings"
)

type (
	// AutoType for the completion
	AutoType rune

	Params struct {
		Command  string   // $1
		Word     string   // $2
		PrevWord string   // $3
		Line     string   // COMP_LINE
		Point    int      // COMP_POINT
		Key      string   // COMP_KEY
		Type     AutoType // COMP_TYPE
	}
)

const (
	TypeNormal AutoType = '\t'
	//TypePartial for listing completions after successive tabs
	TypePartial = '?'
	// TypeList for listing alternatives on partial word completion
	TypeList = '!'
	// TypeListUnmodified to list completions if the word is not unmodified
	TypeListUnmodified = '@'
	//TypeMenu for menu completion
	TypeMenu = '%'
	//TypeNone undefined
	TypeNone = 0
)

//PrepareAutoComplete expects full arg list (from os.Args) and all environment variables (from os.Environ())
func Prepare(args []string, envRaw []string) *Params {
	env := sliceToMap(envRaw)
	var ac Params
	if len(args) != 4 {
		log("expected 4 args, got %v", args)
		return nil
	}
	ac.Command = args[1]
	ac.Word = args[2]
	ac.PrevWord = args[3]
	var ok bool
	//"COMP_LINE", "COMP_TYPE", "COMP_KEY", "COMP_POINT"
	if ac.Line, ok = env["COMP_LINE"]; !ok {
		log("COMP_LINE missing")
		return nil
	}
	if ac.Key, ok = env["COMP_KEY"]; !ok {
		log("COMP_KEY missing")
		return nil
	}
	if sint, ok := env["COMP_TYPE"]; !ok {
		log("COMP_TYPE missing")
		return nil
	} else {
		i, err := strconv.Atoi(sint)
		if err != nil {
			log("COMP_TYPE not an integer")
			return nil
		}
		ac.Type = AutoType(i)
	}
	if sint, ok := env["COMP_POINT"]; !ok {
		log("COMP_POINT missing")
		return nil
	} else {
		i, err := strconv.Atoi(sint)
		if err != nil {
			log("COMP_POINT not an integer")
			return nil
		}
		ac.Point = i
	}

	return &ac
}

func sliceToMap(envRaw []string) map[string]string {
	env := make(map[string]string)
	for _, e := range envRaw {
		ci := strings.Index(e, "=")
		env[e[:ci]] = e[ci+1:]
	}
	return env
}

// LogHook allows callers to intercept the process
var LogHook func(format string, args ...interface{})

func log(format string, args ...interface{}) {
	if LogHook != nil {
		LogHook(format, args...)
	}
}
