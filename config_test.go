package shgf

import "testing"

type confTest struct {
	host  string
	port  int
	valid bool
}

type checkTest struct {
	conf  Config
	valid bool
}

func TestDefaultConf(t *testing.T) {
	var tests = []confTest{
		{"127.0.0.1", 8080, true},
		{"1", 8080, false},
		{"127.0.0.1", 0, false},
		{"wrong", 100000000, false},
	}

	for _, test := range tests {
		if _, err := DefaultConf(test.host, test.port); test.valid && err != nil {
			t.Errorf("expected nil, got %s", err)
		} else if !test.valid && err == nil {
			t.Error("expected error, got nil")
		}
	}
}

func TestConfig_check(t *testing.T) {
	var tests = []checkTest{
		{Config{Hostname: "127.0.0.1", Port: 8080}, true},
		{Config{Hostname: "1", Port: 8080}, false},
		{Config{Hostname: "127.0.0.1", Port: 0}, false},
		{Config{Hostname: "wrong", Port: 100000000}, false},
	}

	for _, test := range tests {
		if err := test.conf.check(); test.valid && err != nil {
			t.Errorf("expected nil, got %s", err)
		} else if !test.valid && err == nil {
			t.Error("expected error, got nil")
		}
	}
}
