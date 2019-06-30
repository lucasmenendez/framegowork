package shgf

import (
	"errors"
	"io/ioutil"
	"mime"
	"os"
	"path/filepath"
	"regexp"
)

const defaultHeader = "text/plain"
const defaultFile = "index.html"

// StaticFolder struct contains the root path for listen requests and has own
// functions to serve static files for this root path.
type StaticFolder struct {
	root string
}

// NewStaticFolder function initializes the folder by the root path provided. If
// the path is not absolute, the function gets the current path and join with it.
func NewStaticFolder(path string) (r *StaticFolder, err error) {
	r = &StaticFolder{path}
	if !filepath.IsAbs(r.root) {
		var current string
		if current, err = os.Getwd(); err != nil {
			return
		}

		r.root = filepath.Join(current, path)
	}

	_, err = os.Stat(r.root)
	return
}

// composePath function joins the path provided with the current StaticFolder
// root path. Then gets the metadata of the resulting path and checks if it is
// safe. If something was wrong, returns an error, else returns the full path.
func (r *StaticFolder) composePath(path string) (string, error) {
	var res = filepath.Join(r.root, path)

	if info, err := os.Stat(res); err != nil {
		return "", err
	} else if info.IsDir() {
		res = filepath.Join(res, defaultFile)
	}

	var outOfScope = regexp.MustCompile(`^\.\/|\.\.`)
	if outOfScope.MatchString(res) {
		return "", errors.New("permission denied")
	}

	return res, nil
}

// Serve function composes a response with the file read, from the filename
// provided. The function calls composePath to get the file path, then gets the
// MIME Type of the file and sets it as response header. Then reads the file and
// writes its content into the response body.
func (r *StaticFolder) Serve(file string) *Response {
	var err error
	var filename string
	if filename, err = r.composePath(file); err != nil {
		if os.IsNotExist(err) {
			return NotFound(err)
		} else if os.IsPermission(err) {
			return Forbidden(err)
		}
		return InternalServerErr(err)
	}

	var res = &Response{Status: 200}
	if header := mime.TypeByExtension(filepath.Ext(filename)); header == "" {
		res.Header = map[string][]string{"Content-type": {defaultHeader}}
	} else {
		res.Header = map[string][]string{"Content-type": {header}}
	}

	if res.Body, err = ioutil.ReadFile(filename); err != nil {
		return InternalServerErr(err)
	}

	return res
}

// StaticFile function splits the path provided into root and filename. With the
// resulting root creates a StaticFolder on-the-fly and returns the response
// composed by the StaticFolder with the resulting filename. If something was
// wrong, catch the resulting error and returns a response according to it.
func StaticFile(path string) *Response {
	var root = filepath.Dir(path)
	var file = filepath.Base(path)

	var err error
	var r *StaticFolder
	if r, err = NewStaticFolder(root); err != nil {
		if os.IsNotExist(err) {
			return NotFound(err)
		} else if os.IsPermission(err) {
			return Forbidden(err)
		}
		return InternalServerErr(err)
	}

	return r.Serve(file)
}
