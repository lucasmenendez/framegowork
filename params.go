package frameworkgo

import "strings"

const routePattern string = "/"

type Params map[string]string

func (params Params) Get(key string) (string, bool) {
	value, ok := params[key]
	return value, ok
}

func cleanComponents(raw_components, pattern string) []string {
	var components []string = strings.Split(raw_components, pattern)

	var cleaned []string
	for _, c := range components {
		if strings.TrimSpace(c) != "" {
			cleaned = append(cleaned, c)
		}
	}

	return cleaned
}

//Extract url params and check if route match with path
func ParseParams(route Route, path string) (bool, Params) {
	var attrs []string
	var params Params = make(Params)

	var routeComponents []string = cleanComponents(route.path, routePattern)
	var pathComponents []string = cleanComponents(path, routePattern)

	if len(routeComponents) == len(pathComponents) {
		for _, s := range routeComponents {
			if len(s) > 0 && string(s[0]) == string(":") {
				attrs = append(attrs, s[1:])
			}
		}

		if route.rgx.MatchString(path) {
			var values []string = route.rgx.FindStringSubmatch(path)[1:]
			for i, v := range values {
				params[attrs[i]] = v
			}

			return true, params
		}
	}
	return false, params
}
