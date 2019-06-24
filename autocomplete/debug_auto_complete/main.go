package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/NearlyUnique/capi/autocomplete"
)

func main() {
	if len(os.Args) != 2 {
		_, _ = fmt.Fprintf(
			os.Stderr,
			"requests exactly 1 argument, the command line you want ti simulate for the complete command\n")
		os.Exit(1)
	}
	autocomplete.EnableLogging(os.Stdout)
	parts := strings.Split(os.Args[1], "|")
	if len(parts) == 1 {
		parts = append(parts, "")
	}
	p := autocomplete.Mock(parts[0], parts[1])
	if p.Word == "" {
		p.Word = `""`
	}

	fmt.Printf(`
# Required env vars
export COMP_LINE=%v
export COMP_POINT=%v
export COMP_KEY=%v
export COMP_TYPE=%v

# IDE env settings
# COMP_LINE=%[1]v;COMP_POINT=%[2]v;COMP_KEY=%[3]v;COMP_TYPE=%[4]v;

# Command Word PreWord
%s %s %s
`,
		p.Line,
		p.Point,
		p.Key,
		p.Type,
		p.Command, p.Word, p.PrevWord)
}
