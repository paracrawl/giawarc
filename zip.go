package giawarc

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"strconv"
	"unicode/utf8"

	"github.com/ulikunitz/xz"
	xzreader "github.com/xi2/xz"
)

type ZipWriter struct {
	fp          *os.File
	compression string
}

type zipWire struct {
	URI         string
	ContentType string
	Lang        string
}

func (zw ZipWriter) Close() (err error) {
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

func (zw ZipWriter) WriteHashes(hashes []uint32) (n int, err error) {
	var buf bytes.Buffer
	var z io.WriteCloser
	if zw.compression == "xz" {
		z, err = xz.NewWriter(&buf)
	} else {
		z = gzip.NewWriter(&buf)
	}
	for _, hash := range hashes{
		fmt.Fprintf(z, "%d\n", hash)
	}

	err = z.Close()
	if err != nil {
		return
	}

	n, err = zw.Write(buf.Bytes())
	return
}

func (zw ZipWriter) WriteText(page *TextRecord) (n int, err error) {
	var buf bytes.Buffer
	var z io.WriteCloser
	if zw.compression == "xz" {
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

// GZ OR XZ READER

type multiStreamReader interface {
	io.Reader
	Multistream(bool)
	Reset(io.Reader) error

}

type GzOrXzReader struct {
	reader	multiStreamReader
	buf	io.Reader
	fp	*os.File
	compression	string
}

func NewGzOrXzReader(compressionType string, filename string) (z GzOrXzReader, err error){
	fp, err := os.Open(filename)
	if err != nil {
		return
	}
	buf := bufio.NewReader(fp)

	var msr multiStreamReader

	if compressionType == "xz" {
		msr, err= xzreader.NewReader(buf, 0)
	} else if compressionType == "gz" {
		msr, err = gzip.NewReader(buf)
	}

	if err != nil{
		return
	}

	z = GzOrXzReader{reader: msr, buf: buf, fp: fp, compression: compressionType}

	return
}

func (z GzOrXzReader) Close() {
	z.fp.Close()
}

func (z GzOrXzReader) Multistream (ms bool) {
	z.reader.Multistream(ms)
}

func (z GzOrXzReader) Reset () (error){
	var err error
	if z.compression == "xz" {
		err = z.reader.Reset(nil)
	} else if z.compression == "gz" {
		err = z.reader.Reset(z.buf)
	}
	return err
}

func (z GzOrXzReader) GetReader () (io.Reader){
	return z.reader
}

func (z GzOrXzReader) ReadHashes() ([]uint32) {
	var hashes []uint32
	reader := bufio.NewReader(z.reader)
	var line string
	var err error
	for {
		line, err = reader.ReadString('\n')
		if err != nil{
			break
		}
		hash, err := strconv.ParseInt(strings.TrimSpace(line), 10, 64)
		if err != nil{
			continue
		}
		hashes = append(hashes, uint32(hash))
	}
	return hashes
}

func (z GzOrXzReader) ReadText() (page *TextRecord, err error) {
	// TODO make this more efficient, it copyies around data too much

	var buf bytes.Buffer
	buf.WriteString("HTTP/1.0 200 OK\n")
	if _, err = buf.ReadFrom(z.reader); err != nil {
		return
	}

	resp, err := http.ReadResponse(bufio.NewReader(&buf), nil)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var t TextRecord
	// t.URI = resp.Header.Get("Content-Location")
	// t.ContentType = resp.Header.Get("Content-Type")
	// t.Lang = resp.Header.Get("Content-Language")
	// t.Date = resp.Header.Get("Date")
	// t.RecordId = resp.Header.Get("X-WARC-Record-ID")
	// t.Source = resp.Header.Get("X-WARC-Filename")
	// t.Text = string(body)

	t.URI = strings.Map(fixUtf, resp.Header.Get("Content-Location"))
	t.ContentType = strings.Map(fixUtf, resp.Header.Get("Content-Type"))
	t.Lang = strings.Map(fixUtf, resp.Header.Get("Content-Language"))
	t.Date = strings.Map(fixUtf, resp.Header.Get("Date"))
	t.RecordId = strings.Map(fixUtf, resp.Header.Get("X-WARC-Record-ID"))
	t.Source = strings.Map(fixUtf, resp.Header.Get("X-WARC-Filename"))
	t.Text = strings.Map(fixUtf, string(body))

	page = &t
	return
}

func fixUtf(r rune) rune {
	if r == utf8.RuneError {
		return -1
	}
	return r
}
