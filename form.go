package shgf

import (
	"net/http"
)

type Form struct {
	keys   []string
	fields map[string]interface{}
}

func parseForm(req *http.Request) (*Form, error) {
	return &Form{}, nil
}

func (f *Form) Exists(key string) bool {
	return false
}

func (f *Form) Get(key string) (interface{}, error) {
	return 0, nil
}

func (f *Form) GetAll() map[string]interface{} {
	return f.fields
}
