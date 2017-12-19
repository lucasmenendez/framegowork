package frameworkgo

import (
	"strings"
	"strconv"
)

const routePattern string = "/"

type Params map[string]string

// Return attribute casted to int from Params if exists
func (p Params) GetInt(key string) (res int, exists bool) {
	var raw string
	if raw, exists = p[key]; exists {
		if res, e := strconv.Atoi(raw); e == nil {
			return res, true
		}
	}

	return res, false
}

// Return attribute casted to float from Params if exists
func (p Params) GetFloat(key string) (res float64, exists bool) {
	var raw string
	if raw, exists = p[key]; exists {
		if res, e := strconv.ParseFloat(raw, 64); e == nil {
			return res, true
		}
	}

	return res, false
}

// Return attribute casted to bool from Params if exists
func (p Params) GetBool(key string) (res bool, exists bool) {
	var raw string
	if raw, exists = p[key]; exists {
		if res, e := strconv.ParseBool(raw); e == nil {
			return res, true
		}
	}

	return res, false
}

// Return raw attribute from Params if exists
func (p Params) GetRaw(key string) (string, bool) {
	value, ok := p[key]
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
