package builder

import (
	"encoding/json"
)

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

type StringOrList []string

func (s *StringOrList) UnmarshalJSON(data []byte) error {
	l := len(data)
	if l >= 2 && data[0] == '"' && data[l-1] == '"' {
		x := []string{string(data[1 : l-1])}
		xs := StringOrList(x)
		*s = xs
		return nil
	}
	list := []string(*s)
	err := json.Unmarshal(data, &list)
	*s = list
	return err
}

func (s StringOrList) MarshalJSON() ([]byte, error) {
	if s == nil {
		return nil, nil
	}
	if len(s) == 1 {
		return json.Marshal([]string(s)[0])
	}
	return json.Marshal([]string(s))
}
