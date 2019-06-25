package giawarc

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"os"
)

type ZipWriter struct {
	*os.File
}

type zipWire struct {
	URI string
	ContentType string
	Lang string
}

func NewZipWriter(out string) (z ZipWriter, err error) {
	fp, err := os.OpenFile(out, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return
	}
	z = ZipWriter{fp}
	return
}

func (zw ZipWriter) WriteText(page *TextRecord) (err error) {
	var buf bytes.Buffer
	z := gzip.NewWriter(&buf)

	z.Name = page.RecordId

	md := []string{page.URI, page.ContentType, page.Lang}
	var meta bytes.Buffer
	enc := gob.NewEncoder(&meta)
	err = enc.Encode(md)
	if err != nil {
		return
	}
	z.Extra = meta.Bytes()

	_, err = z.Write([]byte(page.Text))
	if err != nil {
		return
	}

	err = z.Close()
	if err != nil {
		return
	}

	_, err = zw.Write(buf.Bytes())
	return
}
