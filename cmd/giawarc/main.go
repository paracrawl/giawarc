package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
	"io"
	"bufio"
	"io/ioutil"
	"github.com/paracrawl/giawarc"
)

var outdir string
var outform string
var outputHash string
var inputHash string
var langId string

func init() {
	flag.StringVar(&outdir, "o", ".", "Output location")
	flag.StringVar(&outform, "f", "gzip", "Output format")
	flag.StringVar(&outputHash, "output_hash", "", "Output path Murmur Hash of plain text")
	flag.StringVar(&inputHash, "input_hash", "", "Input path for previous Bitextor Murmur Hash plain texts file")
	flag.StringVar(&langId, "l", "cld2", "Model for language detection: cld2 or cld3")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] WARCFile\nFlags:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(),
`Formats:
  bitextor
        Output format compatible with bitextor (circa June 2019)
  gzip
        Concatenated gzipped pages 
  gzlang
		Concatenated gzipped pages split by language
  xzlang 
		Concatenated xzipped pages split by language
`)
	}
}

func PreProcessFile(filename string) (proc *giawarc.WARCPreProcessor, err error) {
	var f io.ReadCloser
	if filename == "-" {
		reader := bufio.NewReader(os.Stdin)
		f = ioutil.NopCloser(reader)
	} else {
		f, err = os.Open(filename)
	}
	if err != nil {
		return
	}
	defer f.Close()

	var tw giawarc.TextWriter
	if outform == "bitextor" {
		tw, err = giawarc.NewBitextorWriter(outdir, true)
		if err != nil {
			return
		}
	} else if outform == "gzip" {
		tw, err = giawarc.NewZipWriter(outdir, "gz")
		if err != nil {
			return
		}
	} else if outform == "gzlang" {
		m := func(o string) (giawarc.TextWriter, error) { return giawarc.NewZipWriter(o, "gz") }
		tw, err = giawarc.NewLangWriter(outdir, m)
		if err != nil {
			return
		}
	} else if outform == "xzlang" {
		m := func(o string) (giawarc.TextWriter, error) { return giawarc.NewZipWriter(o, "xz") }
		tw, err = giawarc.NewLangWriter(outdir, m)
		if err != nil {
			return
		}
	} else {
		fmt.Fprintf(flag.CommandLine.Output(), "Unknown output format %s\n", outform)
		os.Exit(1)
	}
	defer tw.Close()

	var inputReader giawarc.GzOrXzReader
	if inputHash != "" {
		inputReader, err = giawarc.NewGzOrXzReader("xz", inputHash)
		defer inputReader.Close()
		if err != nil{
			fmt.Println("Error while reading input hashes")
			os.Exit(1)
		}
	}

	var outputWriter giawarc.ZipWriter
	if outputHash != "" {
		outputWriter, err = giawarc.NewZipWriter(outputHash, "xz")
		defer outputWriter.Close()
		if err != nil{
			fmt.Println("Error while opening output hashes file")
			os.Exit(1)
		}
	}

	proc, err = giawarc.NewWARCPreProcessor(f, tw, inputReader, outputWriter, langId)
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

	start := time.Now()
	proc, err := PreProcessFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	end := time.Now()

	elapsed := end.Sub(start)

	log.Printf("total records: %v\n", proc.TotalRecords)
	log.Printf("text records: %v\n", proc.TextRecords)
	log.Printf("lang records: %v\n", proc.LangRecords)
	log.Printf("total bytes: %v\n", proc.TotalBytes)
	log.Printf("text bytes: %v\n", proc.TextBytes)
	log.Printf("lang bytes: %v\n", proc.LangBytes)
	log.Printf("elapsed time: %v\n", elapsed)
}
