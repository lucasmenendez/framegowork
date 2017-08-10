package frameworkgo

import (
	"net/http"
)

const defaultMemory = 32 << 20

type Context struct {
	Path string
	Response http.ResponseWriter
	Request *http.Request
	Handler Handler
	Params map[string]string
}

func NewContext(p string, w http.ResponseWriter, r *http.Request) Context {
	return Context{Path: p, Response: w, Request: r}
}

func (c Context) ParsePostForm() (map[string][]string, error) {
	if err := c.Request.ParseForm(); err != nil {
		return nil, err
	}

	return c.Request.PostForm, nil
}

func (c Context) ParseMultiPartForm() (map[string][]string, error) {
	if err := c.Request.ParseMultipartForm(defaultMemory); err != nil {
		return nil, err
	}

	return c.Request.PostForm, nil
}