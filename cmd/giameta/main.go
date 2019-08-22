package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"github.com/paracrawl/giawarc"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] WARCFile\nFlags:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(),
`Output some metadata about each entry in the WARC file`)
	}
}

func processFile(filename string) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	proc, err := giawarc.NewWARCMetaProcessor(f, path.Base(filename))
	if err != nil {
		return
	}

	proc.Process()

	return
}


func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	filename := flag.Arg(0)

	if err := processFile(filename); err != nil {
		log.Fatal(err)
	}
}
