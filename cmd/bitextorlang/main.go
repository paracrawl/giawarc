package main

import (
	"bufio"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/paracrawl/giawarc"
	"github.com/xi2/xz"
)

var output string
var nrec int
var format string

func init() {
	flag.StringVar(&output, "o", ".", "Output location")
	flag.IntVar(&nrec, "n", -1, "Number of records")
	flag.StringVar(&format, "f", "gz", "Format of input lang files (gz/xz)")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] LangFolder\nFlags:\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func ProcessRecord(reader io.Reader, tw giawarc.TextWriter) (err error) {
	t, err := giawarc.ReadText(reader)
	if err != nil {
		return
	}
	_, err = tw.WriteText(t)
	if err != nil {
		return
	}

	return nil
}

func GZlangToBitextorlang(gzpath string, filename string) (err error) {
	fp, err := os.Open(gzpath)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer fp.Close()

	path := filepath.Join(output, filename)
	// fmt.Println(path)
	os.MkdirAll(path, os.ModePerm)
	tw, err := giawarc.NewBitextorWriter(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer tw.Close()

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
			err = ProcessRecord(zz, tw)
			if err != nil {
				log.Fatal(err)
				return
			}
			err = zz.Reset(buf)
		} else if format == "xz" {
			xx.Multistream(false)
			err = ProcessRecord(xx, tw)
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

	return nil
}

func main() {

	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	gzfolder := flag.Arg(0)

	files, err := ioutil.ReadDir(gzfolder)
	if err != nil {
		return
	}
	for _, f := range files {
		path := filepath.Join(gzfolder, f.Name())
		// fmt.Println("Processing ", path)
		GZlangToBitextorlang(path, f.Name())
	}
}
