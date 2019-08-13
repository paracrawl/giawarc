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
var compression string

func init() {
	flag.StringVar(&field, "o", "uri", "Output field")
	flag.StringVar(&compression, "c", "gz", "Compression format (gz/xz)")
	flag.BoolVar(&b64, "b", false, "Base64 encode output")
	flag.IntVar(&nrec, "n", -1, "Number of records")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] GZFile\nFlags:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "Fields: id, offset, uri, mime, lang, date,  text\n")
	}
}

type MultiStreamReader interface {
	io.Reader
	Multistream(bool)
	Reset(io.Reader) error
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

	var z MultiStreamReader

	if compression == "xz" {
		z, err = xz.NewReader(buf, 0)
	} else if compression == "gz" {
		z, err = gzip.NewReader(buf)
	} else {
		fmt.Println("Unknown compression type: ", compression)
		os.Exit(1)
	}

	if err != nil {
		log.Fatal(err)
		return
	}

	// XXX the xz reader does not implement Close. Since we're
	// operating in reading mode, and the file itself gets
	// closed, this probably doesn't matter too much.
	// defer z.Close()

	for i := 0; i < nrec || nrec == -1; i++ {
		if compression == "gz" {
			z.Multistream(false)
			err = ProcessRecord(z)
			if err != nil {
				log.Fatal(err)
				return
			}
			err = z.Reset(buf)
		}
		if err == io.EOF {
			break
		}
	}
}
