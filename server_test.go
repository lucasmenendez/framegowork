package shgf

import "testing"

type testServer struct {
	address string
	host    string
	port    int
	fail    bool
}

type testAddRoute struct {
	route *route
	fail  bool
}

func TestInitBase(t *testing.T) {
	var tests = []testServer{
		{"", "127.0.0.1", 2020, false},
		{"", "127.0", 2020, true},
		{"", "127.0.0.1", 2020000, true},
		{"", "8.8.8.8", 0, true},
	}

	for _, ti := range tests {
		if _, e := initServer(ti.host, ti.port, false); !ti.fail && e != nil {
			t.Errorf("expected nil, got %s", e)
		} else if ti.fail && e == nil {
			t.Error("expected error, got nil")
		}
	}
}
