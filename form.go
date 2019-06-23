package shgf

import (
	"net/http"
	"strings"
)

const formEncodedContentType = "application/x-www-form-urlencoded"
const formDataContentType = "multipart/form-data"

// Form struct contains parse information from the request body and allows the
// developer to get whole form data or get by key.
type Form struct {
	keys   []string
	fields map[string]interface{}
}

// parseForm function gets a form from the body of the request provided. The
// function decodes the form based on request Content-Type header and generates
// a Form struct with the decoded data. If everything is okey returns a pointer
// of Form, else returns an error.
func parseForm(req *http.Request) (*Form, error) {
	var err error
	var contentType = req.Header.Get("Content-Type")
	if strings.Contains(contentType, formEncodedContentType) {
		err = req.ParseForm()
	} else if strings.Contains(contentType, formDataContentType) {
		err = req.ParseMultipartForm(0)
	}

	if err != nil {
		return &Form{}, err
	}

	var f = Form{
		keys:   []string{},
		fields: make(map[string]interface{}, len(req.Form)),
	}

	for k, v := range req.Form {
		if k != "" {
			f.keys = append(f.keys, k)
			f.fields[k] = v
		}
	}

	return &f, nil
}

// Exists function returns if provided key is currently in the parsed form.
func (f *Form) Exists(key string) bool {
	for _, k := range f.keys {
		if k == key {
			return true
		}
	}

	return false
}

// Get function returns a single form key, if it exists, by the key provided.
func (f *Form) Get(key string) (interface{}, error) {
	if !f.Exists(key) {
		return 0, NewServerErr("form key does not exists")
	}

	return f.fields[key], nil
}

// GetAll function returns whole form into a map whith string keys.
func (f *Form) GetAll() map[string]interface{} {
	return f.fields
}
