package tools

import (
	"encoding/base64"
	"strings"
)

func BasicAuth(header string) (string, string, bool) {
	var username, password string

	ok := false
	if len(header) > 6 {
		authorization, err := base64.StdEncoding.DecodeString(header[6:])
		ok = err == nil

		if ok {
			data := strings.Split(string(authorization), ":")
			if len(data) < 2 {
				ok = false
			} else {
				username, password = data[0], data[1]
			}
		}
	}

	return username, password, ok
}
