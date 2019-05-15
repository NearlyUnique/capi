package capicomplete

import (
	"strings"

	"github.com/NearlyUnique/capi/autocomplete"
	"github.com/NearlyUnique/capi/builder"
)

func indexOrEmpty(args []string, i int) string {
	if i < 0 || i >= len(args) {
		return ""
	}
	return args[i]
}

func GenerateResponse(ac *autocomplete.Params, apis *builder.APISet) []string {
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
		for _, api := range apis.APIs {
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
		//api, err := apis.FindAPI(indexOrEmpty(ac.Args(), 0))
		//if err != nil {
		//	return []string{"error", err.Error()}
		//}
		//for _, cmd := range api.Commands {
		//	// filter
		//	if strings.HasPrefix(cmd.Name, ac.Word) {
		//		//always store simple answer just in case
		//		filtered = append(filtered, cmd.Name)
		//	}
		//	all = append(all, cmd.Name)
		//}
		//if len(filtered) == 0 {
		//	filtered = all
		//}
	default:
		//// must be looking for args
		//api, err := apis.SelectAPI(indexOrEmpty(ac.Args(), 0))
		//if err != nil {
		//	return []string{"no such api", err.Error()}
		//}
		//cmd, err := api.SelectCommand(indexOrEmpty(ac.Args(), 1))
		//if err != nil {
		//	return []string{"no such command", err.Error()}
		//}
		//
		//for _, param := range cmd.ListParams() {
		//	dashed := "--" + param
		//	if strings.HasPrefix(dashed, ac.Word) {
		//		filtered = append(filtered, dashed)
		//	}
		//}
	}

	return filtered
}
