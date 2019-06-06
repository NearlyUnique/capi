package builder

type CommandBody struct {
	Data []byte
}

func (body *CommandBody) UnmarshalJSON(data []byte) error {
	l := len(data)
	if l >= 2 && data[0] == '"' && data[l-1] == '"' {
		data = data[1 : l-1]
	}
	body.Data = data
	return nil
}

func (body *CommandBody) String() string {
	if body == nil {
		return ""
	}
	return string(body.Data)
}
