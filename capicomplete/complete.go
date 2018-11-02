package capicomplete

import (
	"strings"

	"github.com/NearlyUnique/capi"
	"github.com/NearlyUnique/capi/autocomplete"
)

func GenerateResponse(ac *autocomplete.Params, p *capi.Profile) []string {
	var filtered, all []string
	const (
		indexNone    = -1
		indexAPI     = 1
		indexCommand = 2
	)
	index := ac.WordIndex()
	switch index {
	case indexNone:
		if ac.Word != "" {
			break
		}
		fallthrough
	case indexAPI:
		for _, api := range p.APIs {
			// filter
			if strings.HasPrefix(api.Name, ac.Word) {
				//always store simple answer just in case
				filtered = append(filtered, api.Name)
			}
			all = append(all, api.Name)
		}

		if ac.Word == "" && len(filtered) == 0 {
			filtered = all
		}
	case indexCommand:
		api, err := p.SelectAPI(ac.Args())
		if err != nil {
			return []string{"error", err.Error()}
		}
		for _, cmd := range api.Commands {
			// filter
			if strings.HasPrefix(cmd.Name, ac.Word) {
				//always store simple answer just in case
				filtered = append(filtered, cmd.Name)
			}
			all = append(all, cmd.Name)
		}
		if len(filtered) == 0 {
			filtered = all
		}
	}

	return filtered
}
