package giawarc

import (
	"os"
	"path"
)

type LangWriter struct {
	outdir  string
	maker   func(string) (TextWriter, error)
	writers map[string]TextWriter
}

func NewLangWriter(outdir string, maker func(string) (TextWriter, error)) (tw TextWriter, err error) {
	err = os.MkdirAll(outdir, os.ModePerm)
	if err != nil {
		return
	}
	lw := LangWriter{ outdir: outdir }
	lw.maker   = maker
	lw.writers = make(map[string]TextWriter)
	tw = &lw
	return
}

func (lw *LangWriter) WriteText(page *TextRecord) (err error) {
	w, ok := lw.writers[page.Lang]
	if !ok {
		w, err = lw.maker(path.Join(lw.outdir, page.Lang))
		if err != nil {
			return
		}
		lw.writers[page.Lang] = w
	}

	err = w.WriteText(page)
	return
}

func (lw *LangWriter) Close() (err error) {
	for _, w := range(lw.writers) {
		w.Close()
	}
	return
}
