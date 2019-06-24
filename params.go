package shgf

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Params type envolves a map of string and interface{}, and represents a set
// of path params. Provides its own functions to help to the developer to check
// if param exists or get a param by key safely.
type Params map[string]interface{}

const (
	// pathRgx matches with any string between slashes. Used for validation.
	pathRgx = `(\/.*)+`
	// paramDel matches with slashes. Used for split and join paths.
	paramDel = `\/`
	// paramTemp matches with param format. Used for extract params from the
	// path.
	paramTemp = `<(int|string|float|bool):(.*)>`
	// matcherRgx its a template to generate the matcher to handling routes.
	matcherRgx = `(?P<%s_%s>%s)`
	// paramRgx mathes with the route handled param extracting param type and
	// content at once.
	paramRgx = `(int|string|float|bool)_(.+)`
	// floatRgx matches with float number into a string. Used to generate
	// the matcher to handling routes with float params.
	floatRgx = `[0-9]+\.[0-9]+`
	// intRgx matches with int number into a string. Used to generate
	// the matcher to handling routes with int params.
	intRgx = `[0-9]+`
	// boolRgx matches with bool number into a string. Used to generate
	// the matcher to handling routes with bool params.
	boolRgx = `true|false`
	// stringRgx matches with string number into a string. Used to generate
	// the matcher to handling routes with string params.
	strRgx = `.+`
)

// encodeParamRgx function returns a matcher string filling it with the template
// according to the provicded type. Returns error if the type privided is uknown.
func encodeParamRgx(t, a string) (string, error) {
	switch t {
	case "float":
		return fmt.Sprintf(matcherRgx, t, a, floatRgx), nil
	case "int":
		return fmt.Sprintf(matcherRgx, t, a, intRgx), nil
	case "bool":
		return fmt.Sprintf(matcherRgx, t, a, boolRgx), nil
	case "string":
		return fmt.Sprintf(matcherRgx, t, a, strRgx), nil
	}
	return "", NewServerErr("unknown data type")
}

// decodeParamRgx function returns a typed value of provided param based on its
// type (also provided). If the type is unknowed, it returns an error.
func decodeParamRgx(t, v string) (interface{}, error) {
	switch t {
	case "float":
		return strconv.ParseFloat(v, 64)
	case "int":
		return strconv.Atoi(v)
	case "bool":
		return (v == "true"), nil
	}

	return v, nil
}

// encodeParams function returns a regex generated to match with requests path
// defined by path template provided. This template contains named and typed
// variables. encodeParams function encodes each param type and composes a regex
// to match paths with that format. If the template has wrong format, it returns
// an error.
func encodeParams(p string) (*regexp.Regexp, error) {
	var (
		params   []string
		splitter = regexp.MustCompile(paramDel)
		reader   = regexp.MustCompile(paramTemp)
	)

	for _, i := range splitter.Split(p, -1) {
		var t, a string
		if res := reader.FindStringSubmatch(i); len(res) != 3 {
			params = append(params, i)
			continue
		} else {
			t, a = res[1], res[2]
		}

		var err error
		var item string
		if item, err = encodeParamRgx(t, a); err != nil {
			return nil, err
		}

		params = append(params, item)
	}

	var matcherTemp = strings.Join(params, paramDel)
	return regexp.MustCompile(matcherTemp), nil
}

// decodeParams function returns a set of params into a Params struct filled
// with params extracted from the provided path. decodeParams parses each param
// to get its key and decode its value, then store it into a Params struct. If
// any param does not match or something is wrong, it returns an error.
func decodeParams(p string, m *regexp.Regexp) (Params, error) {
	var (
		err      error
		res      = Params{}
		splitter = regexp.MustCompile(paramRgx)
		keys     = m.SubexpNames()
		values   = m.FindStringSubmatch(p)
	)

	keys, values = keys[1:], values[1:]
	for i, k := range keys {
		if k == "" || values[i] == "" {
			return nil, NewServerErr("mismatch labeled params with values provided")
		}

		var metadata = splitter.FindStringSubmatch(k)
		var t, a, v = metadata[1], metadata[2], values[i]
		if res[a], err = decodeParamRgx(t, v); err != nil {
			return nil, err
		}
	}

	return res, nil
}

// Exists function returns if provided key is currently in the parsed params.
func (p Params) Exists(key string) bool {
	_, exists := p[key]
	return exists
}

// Get function returns a single param key, if it exists, by the key provided.
func (p Params) Get(key string) (interface{}, error) {
	if !p.Exists(key) {
		return 0, NewServerErr("form key does not exists")
	}

	return p[key], nil
}
