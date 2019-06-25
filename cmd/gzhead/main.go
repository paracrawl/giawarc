package main

import (
	"bufio"
	"compress/gzip"
	"flag"
	"fmt"
	"log"
	"io"
	"os"
)

var lines int

func init() {
	flag.IntVar(&lines, "n", 10, "Number of Lines")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] GZFile\nFlags:\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	var err error

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

	zr, err := gzip.NewReader(buf)
	if err != nil {
		log.Fatal(err)
	}
	defer zr.Close()

	defer os.Stdout.Close()

	zw := gzip.NewWriter(os.Stdout)

	for i := 0; i < lines; i++ {
		zr.Multistream(false)

		if _, err = io.Copy(zw, zr); err != nil {
			log.Fatal(err)
		}

		if err = zw.Close(); err != nil {
			log.Fatal(err)
		}
		zw.Reset(os.Stdout)

		if err = zr.Reset(buf); err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
	}
}
