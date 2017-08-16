package frameworkgo

import (
	"log"
	"net/http"
	"encoding/json"
)

const defaultMemory = 32 << 20 //32 Mb

type Form map[string]string

func (f Form) Get(key string) (string, bool) {
	val, ok := f[key]
	return val, ok
}

type Context struct {
	Path     string
	response http.ResponseWriter
	request  *http.Request
	Handler  Handler
	Params   map[string]string
}

func NewContext(p string, w http.ResponseWriter, r *http.Request) Context {
	return Context{Path: p, response: w, request: r}
}

func (c Context) ParseForm() (Form, error) {
	if err := c.request.ParseForm(); err != nil {
		return nil, err
	} else if err := c.request.ParseMultipartForm(defaultMemory); err != nil {
		return nil, err
	}

	var form Form = make(map[string]string, len(c.request.PostForm))
	for k, v := range c.request.PostForm {
		form[k] = v[0]
	}

	return form, nil
}

func (c Context) WriteError(err error, status int) {
	c.response.WriteHeader(status)
	c.response.Write([]byte(err.Error()))
}

func (c Context) WriteErrorMessage(err string, status int) {
	c.response.WriteHeader(status)
	c.response.Write([]byte(err))
}

func (c Context) PlainWrite(content []byte, status int) {
	c.response.WriteHeader(status)
	c.response.Write(content)
}

func (c Context) JsonWrite(content interface{}, status int)  {
	if b, err := json.Marshal(content); err != nil {
		log.Fatal(err)
	} else {
		c.response.Header().Set("Content-Type", "application/json")
		c.PlainWrite(b, status)
	}
}