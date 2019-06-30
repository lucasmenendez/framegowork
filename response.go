package shgf

import (
	"encoding/json"
	"log"
	"net/http"
)

// HTTPStatus variable contains a map of int status code and its associated
// []byte message.
var HTTPStatus = map[int][]byte{
	200: []byte("OK"),
	201: []byte("Created"),
	202: []byte("Accepted"),
	203: []byte("Non-Authoritative Information"),
	204: []byte("No Content"),
	205: []byte("Reset Content"),
	206: []byte("Partial Content"),
	207: []byte("Multi-Status"),
	208: []byte("Already Reported"),
	226: []byte("IM Used"),
	300: []byte("Multiple Choices"),
	301: []byte("Moved Permanently"),
	302: []byte("Found"),
	303: []byte("See Other"),
	304: []byte("Not Modified"),
	305: []byte("Use Proxy"),
	307: []byte("Temporary Redirect"),
	308: []byte("Permanent Redirect"),
	400: []byte("Bad Request"),
	401: []byte("Unauthorized"),
	402: []byte("Payment Required"),
	403: []byte("Forbidden"),
	404: []byte("Not Found"),
	405: []byte("Method Not Allowed"),
	406: []byte("Not Acceptable"),
	407: []byte("Proxy Authentication Required"),
	408: []byte("Request Timeout"),
	409: []byte("Conflict"),
	410: []byte("Gone"),
	411: []byte("Length Required"),
	412: []byte("Precondition Failed"),
	413: []byte("Payload Too Large"),
	414: []byte("URI Too Long"),
	415: []byte("Unsupported Media Type"),
	416: []byte("Range Not Satisfiable"),
	417: []byte("Expectation Failed"),
	421: []byte("Misdirected Request"),
	422: []byte("Unprocessable Entity"),
	423: []byte("Locked"),
	424: []byte("Failed Dependency"),
	425: []byte("Too Early"),
	426: []byte("Upgrade Required"),
	428: []byte("Precondition Required"),
	429: []byte("Too Many Requests"),
	431: []byte("Request Header Fields Too Large"),
	451: []byte("Unavailable For Legal Reasons"),
	500: []byte("Internal Server Error"),
	501: []byte("Not Implemented"),
	502: []byte("Bad Gateway"),
	503: []byte("Service Unavailable"),
	504: []byte("Gateway Timeout"),
	505: []byte("HTTP Version Not Supported"),
	506: []byte("Variant Also Negotiates"),
	507: []byte("Insufficient Storage"),
	508: []byte("Loop Detected"),
	510: []byte("Not Extended"),
	511: []byte("Network Authentication Required"),
}

// Response struct its a simplified abstraction of a HTTP response with Header,
// Status and Body.
type Response struct {
	Status int
	Header map[string][]string
	Body   []byte
}

// parseBody function casts a candidate to a valid response body. Casts between
// []byte, string, ServerErr and error to []byte.
func parseBody(d interface{}) (b []byte, e error) {
	switch d.(type) {
	case []byte:
		b = d.([]byte)
		break
	case string:
		b = []byte(d.(string))
		break
	case ServerErr:
		b = []byte(d.(ServerErr).Details())
	case error:
		b = []byte(d.(error).Error())
		break
	default:
		e = NewServerErr("response data must be string, []byte or error")
	}
	return
}

// NewResponse function creates a Response by status int code and a (optional)
// candidate to body. The function creates a default internal server response,
// checks if some data is provided as body and parses it, checks if the provided
// status is valid and returns a Response.
func NewResponse(s int, d ...interface{}) (r *Response, e error) {
	var header = map[string][]string{"Content-type": {"text/plain"}}
	r = &Response{Status: 500, Header: header, Body: []byte(HTTPStatus[500])}

	if len(d) > 1 {
		e = NewServerErr("multiple data provided")
		return
	} else if len(d) == 1 {
		var body []byte
		if body, e = parseBody(d[0]); e != nil {
			return
		}

		r = &Response{s, header, body}
		return
	}

	var desc = HTTPStatus[s]
	if len(desc) == 0 {
		e = NewServerErr("status provided not valid")
		return
	}

	r = &Response{s, header, desc}
	return
}

// JSON function parses provided data interface{} to []byte marshaling it and
// sets HTTP header JSON Content-type and parsed data as Body of current
// Response.
func (r *Response) JSON(d interface{}) (e error) {
	r.Header = map[string][]string{"Content-type": {"application/json"}}
	r.Body, e = json.Marshal(d)
	return
}

// BadRequest function creates a error Response with "Bad Request" HTTP status
// code and sets as Body the default status message with a parsed a set of err
// interface{} as ServerErr.
func BadRequest(err ...interface{}) (r *Response) {
	r, _ = NewResponse(400, NewServerErr(string(HTTPStatus[400]), err...))
	return
}

// Forbidden function creates a error Response with "Forbidden" HTTP status
// code and sets as Body the default status message with a parsed a set of err
// interface{} as ServerErr.
func Forbidden(err ...interface{}) (r *Response) {
	r, _ = NewResponse(403, NewServerErr(string(HTTPStatus[403]), err...))
	return
}

// NotFound function creates a error Response with "Not Found" HTTP status
// code and sets as Body the default status message with a parsed a set of err
// interface{} as ServerErr.
func NotFound(err ...interface{}) (r *Response) {
	r, _ = NewResponse(404, NewServerErr(string(HTTPStatus[404]), err...))
	return
}

// MethodNotAllowed function creates a error Response with "Method Not Allowed"
// HTTP status code and sets as Body the default status message with a parsed a
// set of err interface{} as ServerErr.
func MethodNotAllowed(err ...interface{}) (r *Response) {
	r, _ = NewResponse(405, NewServerErr(string(HTTPStatus[405]), err...))
	return
}

// InternalServerErr function creates a error Response with "Internal Server
// Error" HTTP status code and sets as Body the default status message with a
// parsed a set of err interface{} as ServerErr.
func InternalServerErr(err ...interface{}) (r *Response) {
	r, _ = NewResponse(500, NewServerErr(string(HTTPStatus[500]), err...))
	return
}

// submit function writes Response into a http.ResponseWriter provided. It adds
// each Response Header, the Response Status and Body to http.ResponseWriter and
// write a trace of the response into the log if debug mode is enabled.
func (r *Response) submit(w http.ResponseWriter, debug bool) (e error) {
	for k, vl := range r.Header {
		for _, v := range vl {
			w.Header().Add(k, v)
		}
	}

	w.WriteHeader(r.Status)
	if _, e = w.Write(r.Body); e == nil && debug {
		log.Printf("-> [%d] %s", r.Status, HTTPStatus[r.Status])
	}
	return
}
