package builder

func ListParams(cmd *Command) []string {
	const paramNameIndex = 1
	var list []string
	unique := make(map[string]bool)
	add := func(v string) {
		if _, ok := unique[v]; !ok {
			unique[v] = true
			list = append(list, v)
		}
	}
	if cmd == nil {
		panic("cmd == nil")
	}
	if cmd.API == nil {
		panic("cmd.API == nil")
	}

	searchIn := []string{
		cmd.API.BaseURL,
		cmd.Path,
		cmd.Body.String(),
	}
	headers := []map[string][]string{
		cmd.Query,
		cmd.Header,
		cmd.API.Query,
		cmd.API.Header,
	}

	for _, header := range headers {
		for _, values := range header {
			for _, value := range values {
				searchIn = append(searchIn, value)
			}
		}
	}

	for _, v := range searchIn {
		found := rxMustacheParams.FindAllStringSubmatch(v, -1)
		for _, result := range found {
			add(result[paramNameIndex])
		}
	}

	return list
}
