package giawarc

import (
	"compress/gzip"
	"encoding/base64"
	"io"
	"os"
	"path/filepath"
)

type Zip struct {
	zip io.WriteCloser
	fp  io.WriteCloser
}

func NewZippedFile(outdir string, name string) (z Zip, err error) {
	var zz Zip

	path := filepath.Join(outdir, name)
	zz.fp, err = os.Create(path)
	if err != nil {
		return
	}

	gz, _ := gzip.NewWriterLevel(zz.fp, gzip.BestCompression)
	gz.Name = name
	gz.Comment = "Written by giawarc"
	zz.zip = gz

	z = zz

	return
}

func (z Zip) Write(buf []byte) (int, error) {
	return z.zip.Write(buf)
}

func (z Zip) Close() (err error) {
	z.zip.Close()
	return z.fp.Close()
}

func WriteLine(w io.Writer, s string) (err error) {
	if _, err = w.Write([]byte(s)); err != nil {
		return
	}
	_, err = w.Write([]byte("\n"))
	return
}

type BitextorWriter struct {
  mime  io.WriteCloser
	lang  io.WriteCloser
	url   io.WriteCloser
	plain io.WriteCloser
}

func NewBitextorWriter(outdir string) (tw TextWriter, err error) {
	mime, err := NewZippedFile(outdir, "mime.gz")
	if err != nil {
		return
	}

	lang, err := NewZippedFile(outdir, "lang.gz")
	if err != nil {
		return
	}

	url, err := NewZippedFile(outdir, "url.gz")
	if err != nil {
		return
	}

	plain, err := NewZippedFile(outdir, "plain_text.gz")
	if err != nil {
		return
	}

	return BitextorWriter{mime: mime, lang: lang, url: url, plain: plain}, nil
}

func (bw BitextorWriter) WriteText(text *TextRecord) (n int, err error) {
	if err = WriteLine(bw.mime, text.ContentType); err != nil {
		return
	}

	if err = WriteLine(bw.lang, text.Lang); err != nil {
		return
	}

	if err = WriteLine(bw.url, text.URI); err != nil {
		return
	}

	b64 := base64.StdEncoding.EncodeToString([]byte(text.Text))
	if err = WriteLine(bw.plain, b64); err != nil {
		return
	}

	return
}

func (bw BitextorWriter) Close() (err error) {
	bw.mime.Close()
	bw.lang.Close()
	bw.url.Close()
	bw.plain.Close()
	return
}
