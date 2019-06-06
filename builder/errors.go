package builder

import "fmt"

// InvalidOperation error
type InvalidOperation string

func (e InvalidOperation) Error() string {
	return string(e)
}

//NotFound api or command
type NotFound string

func (e NotFound) Error() string {
	return fmt.Sprintf("search for '%s' returned no results", string(e))
}
