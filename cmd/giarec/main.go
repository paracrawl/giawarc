package main

import (
	"bufio"
	"compress/gzip"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/paracrawl/giawarc"
	"github.com/xi2/xz"
)

var field string
var b64 bool
var nrec int
var format string

func init() {
	flag.StringVar(&field, "o", "uri", "Output field")
	flag.StringVar(&format, "f", "gz", "File format")
	flag.BoolVar(&b64, "b", false, "Base64 encode output")
	flag.IntVar(&nrec, "n", -1, "Number of records")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] GZFile\nFlags:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "Fields: id, offset, uri, mime, lang, date,  text\n")
		fmt.Fprintf(flag.CommandLine.Output(), "File format: gz, xz\n")
	}
}

func ProcessRecord(reader io.Reader) (err error) {
	t, err := giawarc.ReadText(reader)
	if err != nil {
		return
	}

	var out string
	switch field {
	case "id":
		out = t.RecordId
	case "uri":
		out = t.URI
	case "mime":
		out = t.ContentType
	case "date":
		out = t.Date
	case "lang":
		out = t.Lang
	case "text":
		out = t.Text
	}
	if b64 {
		out = base64.StdEncoding.EncodeToString([]byte(out))
	}
	fmt.Printf("%s\n", out)

	return nil
}

func main() {

	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	filename := flag.Arg(0)

	fp, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	buf := bufio.NewReader(fp)
	var xx *xz.Reader
	var zz *gzip.Reader
	if format == "xz" {
		xx, err = xz.NewReader(buf, 0)
	} else if format == "gz" {
		zz, err = gzip.NewReader(buf)
		defer zz.Close()
	} else {
		log.Fatal("Unknown format")
		return
	}

	if err != nil {
		log.Fatal(err)
		return
	}

	for i := 0; i < nrec || nrec == -1; i++ {
		if format == "gz" {
			zz.Multistream(false)
			err = ProcessRecord(zz)
			if err != nil {
				log.Fatal(err)
				return
			}
			err = zz.Reset(buf)
		} else if format == "xz" {
			xx.Multistream(false)
			err = ProcessRecord(xx)
			if err != nil {
				log.Fatal(err)
				return
			}
			err = xx.Reset(nil)
		}
		if err == io.EOF {
			break
		}
	}
}
