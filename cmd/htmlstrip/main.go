package main

import (
	"github.com/paracrawl/giawarc"
	"log"
	"io"
	"os"
)

func main() {
	buf, err := giawarc.HtmlToText(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	io.Copy(os.Stdout, buf)
	os.Stdout.Close()
}
