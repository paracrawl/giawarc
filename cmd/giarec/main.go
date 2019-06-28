package main

import (
	"bufio"
	"compress/gzip"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"io"
	"os"
	"github.com/wwaites/giawarc"
)

var field string
var b64 bool
var nrec int

func init() {
	flag.StringVar(&field, "f", "uri", "Output field")
	flag.BoolVar(&b64, "b", false, "Base64 encode output")
	flag.IntVar(&nrec, "n", -1, "Number of records")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] GZFile\nFlags:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "Fields: id, offset, uri, mime, lang, date,  text\n")
	}
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
	z, err := gzip.NewReader(buf)

	if err != nil {
		log.Fatal(err)
	}
	defer z.Close()

	for i := 0; i < nrec; i++ {
		z.Multistream(false)

		t, err := giawarc.ReadText(z)
		if err != nil {
			log.Fatal(err)
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

		err = z.Reset(buf)
		if err == io.EOF {
			break
		}
	}
}
