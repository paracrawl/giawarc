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
)

var output string
var nrec int

func init() {
	flag.StringVar(&output, "o", ".", "Output location")
	flag.IntVar(&nrec, "n", -1, "Number of records")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] GZFolder\nFlags:\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func GZlangToBitextorlang(gzpath string, filename string) (err error) {
	fp, err := os.Open(gzpath)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer fp.Close()

	path := filepath.Join(output, filename)
	fmt.Println(path)
	os.MkdirAll(path, os.ModePerm)
	tw, err := giawarc.NewBitextorWriter(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer tw.Close()

	buf := bufio.NewReader(fp)
	z, err := gzip.NewReader(buf)

	if err != nil {
		log.Fatal(err)
		return
	}
	defer z.Close()

	for i := 0; i < nrec || nrec == -1; i++ {
		z.Multistream(false)
		t, err := giawarc.ReadText(z)
		if err != nil {
			log.Fatal(err)
		}

		_, err = tw.WriteText(t)
		if err != nil {
			log.Fatal(err)
			return err
		}

		err = z.Reset(buf)
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
		fmt.Println("Processing ", path)
		GZlangToBitextorlang(path, f.Name())
	}
}
