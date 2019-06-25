package giawarc

import (
	"testing"
)

type cleanFixture struct {
	ins  string
	outs string
}

var cleanTests = []cleanFixture{
	{"  hello", "hello"},
	{"hello  \tworld", "hello world"},
	{"\n\nhello", "hello"},
	{"hello \nworld", "hello\nworld"},
	{"hello\n world", "hello\nworld"},
	{"hello\r\nworld", "hello\nworld"},
	{"hello\n\n\nworld", "hello\nworld"},
}

func TestCleanSpace(t *testing.T) {
	for _, fix := range cleanTests {
		o := CleanSpaces(fix.ins)
		if o != fix.outs {
			t.Errorf("%#v != %#v", o, fix.outs)
		}
	}
}
