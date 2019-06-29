package shgf

import (
	"reflect"
	"testing"
)

func TestStatic(t *testing.T) {
	type testStatic struct {
		root string
		res  *StaticRoute
		fail bool
	}

	var currentPath = "/Users/lucasmenendez/Workspace/golang/src/github.com/lucasmenendez/shgf"
	var tests = []testStatic{
		{"./", &StaticRoute{root: currentPath}, false},
		{"./foo", &StaticRoute{}, true},
	}

	for _, test := range tests {
		if res, err := Static(test.root); err != nil && !test.fail {
			t.Errorf("expected nil, got %s", err)
		} else if !reflect.DeepEqual(res, test.res) && !test.fail {
			t.Errorf("expected %+v, got %+v", test.res, res)
		}
	}
}

func TestStaticRoute_composePath(t *testing.T) {
	type testComposePath struct {
		route *StaticRoute
		path  string
		res   string
		fail  bool
	}

	route, _ := Static("./")
	var tests = []testComposePath{
		{route, "./", "/Users/lucasmenendez/Workspace/golang/src/github.com/lucasmenendez/shgf", true},
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
