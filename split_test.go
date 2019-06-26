package giawarc

import (
	"testing"
)


type splitFixture struct {
	lang string
	ins  string
	outs string
}

var splitTests = []splitFixture{
	{"en", "1 hello? World",       "1 hello?\nWorld"},
	{"es", "2 hello! ¡World",      "2 hello!\n¡World"},
	{"en", "3 hello? (World)",     "3 hello?\n(World)"},
	{"en", "4 hello... World",     "4 hello...\nWorld"},
	{"es", "5 hello... (¿World?)", "5 hello...\n(¿World?)"},
	{"en", "6 hello.' World",      "6 hello.'\nWorld"},
	{"en", "7 hello. World",       "7 hello.\nWorld"},
	{"zh", "8 东莞租房|阳江分。房管", "8 东莞租房|阳江分。\n房管"},
	{"zh", "9 Hello.大全|百度阳江吧|阳江商友", "9 Hello.\n大全|百度阳江吧|阳江商友"},
	{"en", "10 Good morning Mr. Wu", "10 Good morning Mr. Wu"},
	{"en", "11 I'm No. 1", "11 I'm No. 1"},
	{"en", "12 No.", "12 No."},
	{"en", "13 (Luke 7:1–10). First", "13 (Luke 7:1–10).\nFirst"},
}

func TestSplitSentences(t *testing.T) {
	for _, fix := range splitTests {
		o := SplitSentences(fix.ins, fix.lang)
		if o != fix.outs {
			t.Errorf("%#v != %#v", o, fix.outs)
		}
	}
}
