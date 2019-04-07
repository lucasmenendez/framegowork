package shgf

import (
	"errors"
	"fmt"
)

// unknown constant contains default message to unknown error
const unknown = "unknown error"

// ServerErr struct is a custom error for shgf that contains and simple message
// and original error.
type ServerErr struct {
	msg string
	err error
}

// NewServerErr function creates new ServerErr by message and (optional) error
// provided as arguments. The function assing the message provided to custom
// ServerErr, checks if error exists and casts it to error.
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

// Error function implements Error() function from Error interface.
func (e ServerErr) Error() string {
	return fmt.Sprintf("shgf error: %s", e.msg)
}

// Details function returns merged string between ServerErr message and error.
func (e ServerErr) Details() string {
	if e.err != nil {
		return fmt.Sprintf("shf error: %s\n%s", e.msg, e.err.Error())
	}

	return e.Error()
}
