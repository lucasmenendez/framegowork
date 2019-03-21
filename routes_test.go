package shgf

import (
	"regexp"
	"testing"
)

var routesTests = routes{
	&route{method: "POST", path: "/foo1/bar", matcher: regexp.MustCompile(`\/foo1\/bar`)},
	&route{method: "POST", path: "/foo2/bar", matcher: regexp.MustCompile(`\/foo2\/bar`)},
	&route{method: "POST", path: "/foo3", matcher: regexp.MustCompile(`\/foo3`)},
	&route{method: "POST", path: "/foo4/bar", matcher: regexp.MustCompile(`\/foo4\/bar`)},
}

func TestRoutes_exists(t *testing.T) {
	if !routesTests.exists(route{method: "POST", path: "/foo4/bar"}) {
		t.Error("expected true, got false")
	} else if routesTests.exists(route{method: "GET", path: "/foo1/bar"}) {
		t.Error("expected false, got true")
	} else if routesTests.exists(route{method: "PUT"}) {
		t.Error("expected false, got true")
	}
}

func TestRoutes_findByPath(t *testing.T) {
	if r := routesTests.findByPath("/foo4/bar"); len(r) != 1 {
		t.Error("expected true, got false")
	} else if r := routesTests.findByPath("/foo/bar2"); len(r) != 0 {
		t.Error("expected false, got true")
	} else if r := routesTests.findByPath(""); len(r) != 0 {
		t.Error("expected false, got true")
	}
}

func TestRoutesFindByMatchingPath(t *testing.T) {
	if r := routesTests.findByPath("/foo4/bar"); len(r) != 1 {
		t.Error("expected true, got false")
	} else if r := routesTests.findByPath("/foo/bar2"); len(r) != 0 {
		t.Error("expected false, got true")
	} else if r := routesTests.findByPath(""); len(r) != 0 {
		t.Error("expected false, got true")
	}
}

func TestRoutes_findByMethod(t *testing.T) {
	if r := routesTests.findByMethod("POST"); len(r) != 4 {
		t.Error("expected true, got false")
	} else if r := routesTests.findByMethod("GET"); len(r) != 0 {
		t.Error("expected false, got true")
	} else if r := routesTests.findByMethod(""); len(r) != 0 {
		t.Error("expected false, got true")
	}
}
