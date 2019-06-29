package shgf

import (
	"net/http"
	"os"
	"path/filepath"
)

type StaticRoute struct {
	root string
	fs   http.FileSystem
}

func Static(path string) (r *StaticRoute, err error) {
	var current string
	if current, err = os.Getwd(); err != nil {
		return
	}

	r = &StaticRoute{}
	r.root = filepath.Join(current, path)
	_, err = os.Stat(r.root)
	return
}

func (r *StaticRoute) composePath(path string) (string, error) {
	var res = filepath.Join(r.root, path)

	if _, err := os.Stat(res); os.IsNotExist(err) {
		return "", err
	}
	return res, nil
}
