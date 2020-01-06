package giawarc

import (
	"encoding/base64"
	"errors"
	"io"
	"os"
	"path/filepath"
)

type BiLangWriterList struct {
	mime  io.WriteCloser
	url   io.WriteCloser
	plain io.WriteCloser
}

type BiLangWriter struct {
	outdir string
	ws map[string]BiLangWriterList
}

func NewBiLangWriter(outdir string, writeLang bool) (tw TextWriter, err error) {
	bw := BiLangWriter{outdir: outdir, ws: make(map[string]BiLangWriterList)}
	return &bw, nil
}

func (bw BiLangWriter) WriteText(text *TextRecord) (n int, err error) {
	if len(text.Lang) == 0 {
		err = errors.New("Invalid empty language")
		return
	}

	ws, ok := bw.ws[text.Lang]
	if !ok {
		outsubdir := filepath.Join(bw.outdir, text.Lang)
                err = os.MkdirAll(outsubdir, os.ModePerm)
		if err != nil {
			return
		}

		mime, err := NewZippedFile(outsubdir, "mime.gz")
		if err != nil {
			return 0, err
		}

		url, err := NewZippedFile(outsubdir, "url.gz")
		if err != nil {
			mime.Close()
			return 0, err
		}

		plain, err := NewZippedFile(outsubdir, "plain_text.gz")
		if err != nil {
			mime.Close()
			url.Close()
			return 0, err
		}

		ws = BiLangWriterList{mime: mime, url: url, plain: plain}
		bw.ws[text.Lang] = ws
	}

	if err = WriteLine(ws.mime, text.ContentType); err != nil {
		return
	}

	if err = WriteLine(ws.url, text.URI); err != nil {
		return
	}

	b64 := base64.StdEncoding.EncodeToString([]byte(text.Text))
	if err = WriteLine(ws.plain, b64); err != nil {
		return
	}

	return
}

func (bw BiLangWriter) Close() (err error) {
	for _, ws := range bw.ws {
		ws.mime.Close()
		ws.url.Close()
		ws.plain.Close()
	}
	return
}
