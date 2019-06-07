package capicomplete

import (
	"fmt"
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
		apiName := indexOrEmpty(ac.Args(), 0)
		api, err := apis.FindAPI(apiName)
		if err != nil {
			return []string{"error", err.Error()}
		}
		if len(api) != 1 {
			return []string{"error", fmt.Sprintf("API '%s' ambiguous", apiName)}
		}
		cmds, err := api[0].FindCommand(ac.Word)
		if err != nil {
			return []string{"error", err.Error()}
		}
		for _, cmd := range cmds {
			if strings.HasPrefix(cmd.Name, ac.Word) {
				//always store simple answer just in case
				filtered = append(filtered, cmd.Name)
			}
			all = append(all, cmd.Name)
		}
		if len(filtered) == 0 {
			filtered = all
		}
	default:
		// must be looking for args
		apiName := indexOrEmpty(ac.Args(), 0)
		api, err := apis.FindAPI(apiName)
		if err != nil {
			return []string{"error", err.Error()}
		}
		if len(api) != 1 {
			return []string{"error", fmt.Sprintf("API '%s' ambiguous", apiName)}
		}
		cmdName := indexOrEmpty(ac.Args(), 1)
		cmd, err := api[0].FindCommand(cmdName)
		if err != nil {
			return []string{"error", err.Error()}
		}
		if len(cmd) != 1 {
			return []string{"error", fmt.Sprintf("command '%s' ambiguous", cmdName)}
		}

		for _, param := range builder.ListParams(cmd[0]) {
			dashed := "--" + param
			if strings.HasPrefix(dashed, ac.Word) {
				filtered = append(filtered, dashed)
			}
		}
	}

	return filtered
}
