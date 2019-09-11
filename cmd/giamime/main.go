package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/paracrawl/go-warc/warc"
	"github.com/paracrawl/giawarc"
)

var mimes map[string]int

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] WARCFile\nFlags:\n", os.Args[0])
		flag.PrintDefaults()
	}

	mimes = make(map[string]int)
}

func processRecord(wr *warc.WARCRecord, err error) {
	if err != nil {
		return
	}

	// content type of the WARC record not the payload
	content_type, _ := wr.GetHeader().Get("Content-Type")
	if !strings.Contains(content_type, "application/http") || !strings.Contains(content_type, "response") {
		// nothing to do
		// log.Printf("Ignoring WARC record (not application/http)")
		return
	}

	// get HTTP response out of the WARC file, and parse it
	payload := wr.GetPayload()
	resp, err := http.ReadResponse(bufio.NewReader(payload.GetReader()), nil)
	if err == nil {
		content_type, _ = giawarc.CleanContentType(resp.Header.Get("Content-Type"))
		c := mimes[content_type]
		mimes[content_type] = c + 1
	}
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	filename := flag.Arg(0)

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	wf, err := warc.NewWARCFile(f)
	if err != nil {
		log.Fatal(err)
	}

	wf.GetReader().Iterate(processRecord)

	for m, c := range mimes {
		fmt.Printf("%v\t%v\n", c, m)
	}
}
