package autocomplete

import (
	"fmt"
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

//Prepare expects full arg list (from os.Args) and all environment variables (from os.Environ())
func Prepare(args []string, envRaw []string) *Params {
	const (
		indexCommand  = 1
		indexWord     = 2
		indexPrevWord = 3
	)
	env := sliceToMap(envRaw)
	var ac Params
	if len(args) != 4 {
		log("expected 4 args, got %v", args)
		return nil
	}
	ac.Command = args[indexCommand]
	ac.Word = args[indexWord]
	ac.PrevWord = args[indexPrevWord]
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

//WordIndex for the word under the cursor
func (p *Params) WordIndex() int {
	return logicalPosition(p.Line, p.Point)
}

// Args as would be specified by os.Args
func (p *Params) Args() []string {
	return strings.Split(p.Line, " ")
}

func sliceToMap(envRaw []string) map[string]string {
	env := make(map[string]string)
	for _, e := range envRaw {
		ci := strings.Index(e, "=")
		env[e[:ci]] = e[ci+1:]
	}
	return env
}

func logicalPosition(line string, pos int) int {
	if pos > len(line) {
		return -1
	}
	p := 0
	for i := 0; i < pos; i++ {
		if line[i] == ' ' {
			p++
		}
	}
	return p
}

// LogHook allows callers to intercept the process
var LogHook func(format string, args ...interface{})

func log(format string, args ...interface{}) {
	if LogHook != nil {
		LogHook(format, args...)
	}
}

func (p Params) String() string {
	return fmt.Sprintf("Command=%v,Word=%v,PrevWord=%v,Line=%v,Point=%v,Key=%v,Type=%v",
		p.Command, p.Word, p.PrevWord, p.Line, p.Point, p.Key, p.Type)
}
