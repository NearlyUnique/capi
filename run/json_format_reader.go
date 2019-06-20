package run

import (
	"encoding/json"

	"github.com/NearlyUnique/capi/builder"
)

const (
	JSONFileExtension = ".json"
)

func JSONFormatReader(content []byte) (*builder.APISet, error) {
	var set builder.APISet
	err := json.Unmarshal(content, &set)
	if err != nil {
		return nil, err
	}
	return &set, nil
}
