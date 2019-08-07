package giawarc

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"github.com/ulikunitz/xz"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type ZipWriter struct {
	fp *os.File
	compression string
}

type zipWire struct {
	URI string
	ContentType string
	Lang string
}

func (zw ZipWriter) Close() (err error){
	return zw.fp.Close()
}

func (zw ZipWriter) Write(buf []byte) (n int, err error) {
	return zw.fp.Write(buf)
}

func NewZipWriter(out string, compression string) (z ZipWriter, err error) {
	fp, err := os.OpenFile(out, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return
	}
	z = ZipWriter{fp: fp, compression: compression}
	return
}

func (zw ZipWriter) WriteText(page *TextRecord) (n int, err error) {
	var buf bytes.Buffer
	var z io.WriteCloser
	if zw.compression == "xz"{
		z, err = xz.NewWriter(&buf)
	} else {
		z = gzip.NewWriter(&buf)
	}
	// z.Name = page.RecordId

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


func ReadText(z io.Reader) (page *TextRecord, err error) {
	// TODO make this more efficient, it copyies around data too much

	var buf bytes.Buffer
	buf.WriteString("HTTP/1.0 200 OK\n")
	if _, err = buf.ReadFrom(z); err != nil {
		return
	}
	fmt.Println(buf.String())

	resp, err := http.ReadResponse(bufio.NewReader(&buf), nil)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var t TextRecord
	t.URI  = resp.Header.Get("Content-Location")
	t.ContentType = resp.Header.Get("Content-Type")
	t.Lang = resp.Header.Get("Content-Language")
	t.Date = resp.Header.Get("Date")
	t.RecordId = resp.Header.Get("X-WARC-Record-ID")
	t.Source = resp.Header.Get("X-WARC-Filename")
	t.Text = string(body)


	page = &t
	return
}

