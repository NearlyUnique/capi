package capicomplete

import (
	"fmt"
	"strings"

	"github.com/NearlyUnique/capi"
	"github.com/NearlyUnique/capi/autocomplete"
)

func GenerateResponse(ac *autocomplete.Params, p *capi.Profile) []string {
	var all []string
	single := []string{""}
	hasPrefix := func(api capi.API) bool {
		return strings.HasPrefix(api.Name, ac.Word) ||
			strings.HasPrefix(api.Alias, ac.Word)
	}
	alwaysInclude := func(api capi.API) bool { return true }

	for _, shouldInclude := range []func(api capi.API) bool{hasPrefix, alwaysInclude} {
		for _, api := range p.APIs {
			// filter
			if shouldInclude(api) {
				//always store simple answer just in case
				if len(single[0]) == 0 {
					single[0] = api.Name
				}
				// store
				if api.Alias == "" {
					all = append(all, api.Name)
				} else {
					all = append(all, fmt.Sprintf("(%s) %s", api.Alias, api.Name))
				}
			}
		}
		if len(all) == 1 {
			return single
		} else if len(all) > 1 {
			break
		}
	}

	return all
}
