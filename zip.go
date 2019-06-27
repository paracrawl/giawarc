package giawarc

import (
	"bytes"
	"compress/gzip"
	"fmt"
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

func (zw ZipWriter) WriteText(page *TextRecord) (n int, err error) {
	var buf bytes.Buffer
	z := gzip.NewWriter(&buf)

	z.Name = page.RecordId

	fmt.Fprintf(z, "Content-Location: %s\n", page.URI)
	fmt.Fprintf(z, "Content-Type: %s\n", page.ContentType)
	fmt.Fprintf(z, "Content-Language: %s\n", page.Lang)
	fmt.Fprintf(z, "Content-Length: %d\n", len(page.Text))
	fmt.Fprintf(z, "Date: %s\n", page.Date)
	fmt.Fprintf(z, "X-WARC-Record-ID: <urn:uuid:%s>\n", page.RecordId)
	fmt.Fprintf(z, "X-WARC-Filename: %s\n", page.Source)
	fmt.Fprintf(z, "\n")

	_, err = z.Write([]byte(page.Text))
	if err != nil {
		return
	}
	fmt.Fprintf(z, "\n")

	err = z.Close()
	if err != nil {
		return
	}

	n, err = zw.Write(buf.Bytes())
	return
}
