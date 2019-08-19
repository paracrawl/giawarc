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
var compression string
var format string

type MultiStreamReader interface {
	io.Reader
	Multistream(bool)
	Reset(io.Reader) error
}

func init() {
	flag.StringVar(&output, "o", ".", "Output location")
	flag.IntVar(&nrec, "n", -1, "Number of records")
	flag.StringVar(&compression, "c", "gz", "Compression of input lang files (gz/xz)")
	flag.StringVar(&format, "f", "bitextorlang", "Output format (bitextor/bitextorlang)")
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

func GZlangToBitextorlang(tw giawarc.TextWriter, gzpath string) (err error) {
	fp, err := os.Open(gzpath)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer fp.Close()

	if err != nil {
		log.Fatal(err)
		return
	}

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
		return
	}

	for i := 0; i < nrec || nrec == -1; i++ {
		z.Multistream(false)
		err = ProcessRecord(z, tw)
		if err != nil {
			log.Fatal(err)
			return
		}
		if compression == "gz"{
			err = z.Reset(buf)
		} else {
			err = z.Reset(nil)
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
	var tw giawarc.TextWriter
	// if output format is bitextor, only one writer needed as everything goes into the same directory
	if format=="bitextor" {
		os.MkdirAll(output, os.ModePerm)
		tw, err = giawarc.NewBitextorWriter(output, true) 
	} else if format=="bitextorlang"{
	}else{
		fmt.Println("Unkown output format: ", format)
		os.Exit(1)
	}
	for _, f := range files {
		inputPath := filepath.Join(gzfolder, f.Name())
		// if output format is bitextorlang, create a writer for each lang directory 
		if format=="bitextorlang"{
			outputPath := filepath.Join(output, f.Name())
			os.MkdirAll(outputPath, os.ModePerm)
			tw, err = giawarc.NewBitextorWriter(outputPath, false)
			if err != nil {
				log.Fatal(err)
				return
			}
		}
		// fmt.Println("Processing ", inputPath)
		GZlangToBitextorlang(tw, inputPath)
		if format=="bitextorlang"{
			tw.Close()
		}
	}
	if format == "bitextor"{
		tw.Close()
	}
}
