package shgf

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewStaticFolder(t *testing.T) {
	type testStatic struct {
		root string
		res  *StaticFolder
		fail bool
	}

	var currentPath = "/Users/lucasmenendez/Workspace/golang/src/github.com/lucasmenendez/shgf"
	var tests = []testStatic{
		{"./", &StaticFolder{root: currentPath}, false},
		{"./foo", &StaticFolder{}, true},
	}

	for _, test := range tests {
		if res, err := NewStaticFolder(test.root); err != nil && !test.fail {
			t.Errorf("expected nil, got %s", err)
		} else if !reflect.DeepEqual(res, test.res) && !test.fail {
			t.Errorf("expected %+v, got %+v", test.res, res)
		}
	}
}

func TestStaticFolder_composePath(t *testing.T) {
	type testComposePath struct {
		route *StaticFolder
		path  string
		res   string
		fail  bool
	}

	route, _ := NewStaticFolder("./")
	var tests = []testComposePath{
		{route, "./", "/Users/lucasmenendez/Workspace/golang/src/github.com/lucasmenendez/shgf/index.html", true},
		{route, "config.go", "/Users/lucasmenendez/Workspace/golang/src/github.com/lucasmenendez/shgf/config.go", false},
		{route, "foo", "", true},
	}

	for _, test := range tests {
		if res, err := test.route.composePath(test.path); err != nil && !test.fail {
			t.Errorf("expected nil, got %s", err)
		} else if res != test.res {
			t.Errorf("expected %s, got %s", test.res, res)
		}
	}
}

func TestStaticFolderServe(t *testing.T) {
	var r, _ = NewStaticFolder("./")

	type serveTests struct {
		file    string
		content []byte
		res     *Response
		fail    bool
	}

	var tests = []serveTests{
		{"test.json", []byte(`{"foo": { "bar": 1 }}`), &Response{200, map[string][]string{"Content-type": {"text/plain"}}, []byte(`{"foo": { "bar": 1 }}`)}, true},
		{"test.json", []byte(`{"foo": { "bar": 1 }}`), &Response{200, map[string][]string{"Content-type": {"application/json"}}, []byte(`{"foo": { "bar": 1 }}`)}, false},
	}

	for _, test := range tests {
		var filename = filepath.Join(r.root, test.file)
		ioutil.WriteFile(filename, test.content, os.ModePerm)
		defer os.Remove(filename)

		if res := r.Serve(test.file); !reflect.DeepEqual(res, test.res) && !test.fail {
			t.Errorf("expected %+v, got %+v", test.res, res)
		}
	}
}
