package giawarc

import (
	"compress/gzip"
	"encoding/base64"
	"io"
	"os"
	"path/filepath"
	"github.com/ulikunitz/xz"
)

type Zip struct {
	zip io.WriteCloser
	fp  io.WriteCloser
}

type XZip struct {
	xzip io.WriteCloser
	fp   io.WriteCloser
}

func NewZippedFile(outdir string, name string) (z Zip, err error) {
	var zz Zip

	path := filepath.Join(outdir, name)
	zz.fp, err = os.Create(path)
	// zz.fp, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func NewXZipFile(outdir string, name string) (x XZip, err error) {
	var xx XZip

	path := filepath.Join(outdir, name)
	xx.fp, err = os.Create(path)
	// xx.fp, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return xx, err
	}

	w, err := xz.NewWriter(xx.fp)
	xx.xzip = w
	return xx, err
}

func (z Zip) Write(buf []byte) (int, error) {
	return z.zip.Write(buf)
}

func (x XZip) Write(buf []byte) (int, error) {
	return x.xzip.Write(buf)
}

func (z Zip) Close() (err error) {
	z.zip.Close()
	return z.fp.Close()
}

func (x XZip) Close() (err error) {
	x.xzip.Close()
	return x.fp.Close()
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
	id    io.WriteCloser
}

func NewBitextorWriter(outdir string, writeLang bool) (tw TextWriter, err error) {
	mime, err := NewXZipFile(outdir, "mime.xz")
	if err != nil {
		return
	}

	var lang io.WriteCloser
	lang = nil
	if writeLang {
		lang, err = NewXZipFile(outdir, "lang.xz")
		if err != nil {
			return
		}
	}


	url, err := NewXZipFile(outdir, "url.xz")
	if err != nil {
		return
	}

	plain, err := NewXZipFile(outdir, "plain_text.xz")
	if err != nil {
		return
	}

	id, err := NewXZipFile(outdir, "uuid.xz")
	if err != nil {
		return
	}

	return BitextorWriter{mime: mime, lang: lang, url: url, plain: plain, id: id}, nil
}

func (bw BitextorWriter) WriteText(text *TextRecord) (n int, err error) {
	if err = WriteLine(bw.mime, text.ContentType); err != nil {
		return
	}

	if (bw.lang != nil){
		if err = WriteLine(bw.lang, text.Lang); err != nil {
			return
		}
	}

	if err = WriteLine(bw.url, text.URI); err != nil {
		return
	}

	b64 := base64.StdEncoding.EncodeToString([]byte(text.Text))
	if err = WriteLine(bw.plain, b64); err != nil {
		return
	}

	if err = WriteLine(bw.id, text.RecordId); err != nil {
		return
	}
	return
}

func (bw BitextorWriter) Close() (err error) {
	bw.mime.Close()
	if bw.lang != nil {
		bw.lang.Close()
	}
	bw.url.Close()
	bw.plain.Close()
	bw.id.Close()
	return
}
