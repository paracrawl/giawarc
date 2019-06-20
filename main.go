package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func PreProcessFile(filename string, outdir string) (proc *WARCPreProcessor, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	bw, err := NewBitextorWriter("./")
	if err != nil {
		return
	}
	defer bw.Close()

	proc, err = NewWARCPreProcessor(f, bw)
	if err != nil {
		return
	}

	proc.Process()

	return
}

func main() {
	filename := os.Args[1]

	start := time.Now()
	proc, err := PreProcessFile(filename, "./")
	if err != nil {
		log.Fatal(err)
	}
	end := time.Now()

	elapsed := end.Sub(start)

	fmt.Printf("total records: %v\n", proc.TotalRecords)
	fmt.Printf("text records: %v\n", proc.TextRecords)
	fmt.Printf("lang records: %v\n", proc.LangRecords)
	fmt.Printf("total bytes: %v\n", proc.TotalBytes)
	fmt.Printf("text bytes: %v\n", proc.TextBytes)
	fmt.Printf("lang bytes: %v\n", proc.LangBytes)
	fmt.Printf("elapsed time: %v\n", elapsed)
	fmt.Printf("content types:\n")

	cts := proc.ContentTypeStats()
	for _, s := range cts {
		fmt.Printf("    %v: %0.08f\n", s.ContentType, s.Prevalence)
	}
}
