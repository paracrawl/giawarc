package giawarc

import (
	"testing"
)

type cleanFixture struct {
	ins  string
	outs string
}

var cleanTests = []cleanFixture{
	{"  hello", "hello\n"},
	{"hello  \tworld", "hello world\n"},
	{"\n\nhello", "hello\n"},
	{"hello \nworld", "hello\nworld\n"},
	{"hello\n world", "hello\nworld\n"},
	{"hello\r\nworld", "hello\nworld\n"},
	{"hello\n\n\n\nworld\n\nxyz", "hello\nworld\nxyz\n"},
}

func TestCleanSpace(t *testing.T) {
	for _, fix := range cleanTests {
		o := CleanSpaces(fix.ins)
		if o != fix.outs {
			t.Errorf("%#v != %#v", o, fix.outs)
		}
	}
}
