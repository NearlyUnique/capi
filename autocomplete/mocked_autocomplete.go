package autocomplete

import "strings"

// Mock an input, before and after the cursor position. after may be empty
// expected to be the full cli line including the command as seen in the shell
func Mock(before, after string) Params {
	p := Params{
		Line:  before + after,
		Point: len(before),
		Key:   "9",
		Type:  TypeNormal,
	}
	i := strings.Index(before, " ")
	if i < 0 {
		p.Command = before
	} else {
		p.Command = before[:i]
	}

	const (
		scan = -1
		word = 0
		prev = 1
	)

	state := word
	ilast := len(before)

	i = ilast - 1
loop:
	for ; i >= 0; i-- {
		if before[i] == ' ' {
			switch state {
			case word:
				state = scan
				p.Word = before[i+1 : ilast]
			case prev:
				break loop
			}
			ilast = i
		} else {
			if state == scan {
				state = prev
				ilast = i + 1
			}
		}
	}
	p.PrevWord = before[i+1 : ilast]

	return p
}
