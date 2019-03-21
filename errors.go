package shgf

import (
	"errors"
	"fmt"
)

const unknown = "unknown error"

type ServerErr struct {
	msg string
	err error
}

func NewServerErr(msg string, err ...interface{}) (e ServerErr) {
	e.msg = msg

	if len(err) > 0 {
		switch err[0].(type) {
		case error:
			e.err = err[0].(error)
			break
		case string:
			e.err = errors.New(err[0].(string))
			break
		case []byte:
			e.err = errors.New(string(err[0].([]byte)))
			break
		default:
			e.err = errors.New(unknown)
		}
	}

	return
}

func (e ServerErr) Error() string {
	return fmt.Sprintf("shgf error: %s", e.msg)
}

func (e ServerErr) Details() string {
	return fmt.Sprintf("shf error: %s\n%s", e.msg, e.err.Error())
}
