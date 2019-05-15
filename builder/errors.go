package builder

// InvalidOperation error
type InvalidOperation string

func (e InvalidOperation) Error() string {
	return string(e)
}
