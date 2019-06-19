package autocomplete

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
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

		err error
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

//Parse expects full arg list (from os.Args) and all environment variables (from os.Environ())
func Parse(args []string, envRaw []string) *Params {
	const (
		indexCommand  = 1
		indexWord     = 2
		indexPrevWord = 3
	)
	logArgs(args, envRaw)

	var ac Params
	if len(args) != 4 {
		log("expected 4 args, got %v", args)
		return nil
	}
	env := sliceToMap(envRaw)

	ac.Command = args[indexCommand]
	ac.Word = args[indexWord]
	ac.PrevWord = args[indexPrevWord]

	ac.parseCompLine(env)
	ac.parseCompKey(env)
	ac.parseCompType(env)
	ac.ParseCompPoint(env)

	if ac.err != nil {
		log(ac.err.Error())
		return nil
	}

	return &ac
}

func logArgs(args, envRaw []string) {
	log("args:")
	for i, v := range args {
		log("%d:'%s',", i, v)
	}
	log("\nenv:")
	for _, v := range envRaw {
		if strings.HasPrefix(v, "COMP_") {
			log("%s;", v)
		}
	}
	log("\n")
}

//WordIndex for the word under the cursor
func (p *Params) WordIndex() int {
	return logicalPosition(p.Line, p.Point)
}

// Args only the args to the command, not the command itself, use CliArgs for that
func (p *Params) Args() []string {
	args := p.CliArgs()
	if len(args) == 1 {
		return nil
	}
	return args[1:]
}

// CliArgs as would be specified by os.Args
func (p *Params) CliArgs() []string {
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

func (p *Params) parseCompLine(env map[string]string) {
	if p.err != nil {
		return
	}
	var ok bool
	p.Line, ok = env["COMP_LINE"]
	if !ok {
		p.err = errors.New("COMP_LINE missing")
	}
}

func (p *Params) parseCompKey(env map[string]string) {
	if p.err != nil {
		return
	}
	var ok bool
	p.Key, ok = env["COMP_KEY"]
	if !ok {
		p.err = errors.New("COMP_KEY missing")
	}

}

func (p *Params) parseCompType(env map[string]string) {
	if p.err != nil {
		return
	}
	temp, ok := env["COMP_TYPE"]
	if !ok {
		p.err = errors.New("COMP_TYPE missing")
		return
	}
	i, err := strconv.Atoi(temp)
	if err != nil {
		p.err = errors.New("COMP_TYPE not an integer")
		return
	}
	p.Type = AutoType(i)
}

func (p *Params) ParseCompPoint(env map[string]string) {
	if p.err != nil {
		return
	}
	temp, ok := env["COMP_POINT"]
	if !ok {
		p.err = errors.New("COMP_POINT missing")
		return
	}

	i, err := strconv.Atoi(temp)
	if err != nil {
		p.err = errors.New("COMP_POINT not an integer")
	}
	p.Point = i
}
