package frameworkgo

import "strings"

type Params map[string]string

//Extract url params and check if route match with path
func ParseParams(route Route, path string) (bool, Params) {
	var attrs []string
	var params Params = make(Params)

	var routeComponents []string = strings.Split(route.path, "/")
	var pathComponents []string = strings.Split(path, "/")

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
		} else {
			return false, params
		}
	} else {
		return false, params
	}
}
