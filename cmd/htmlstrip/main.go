package main

import (
	"flag"
	"github.com/paracrawl/giawarc"
	"log"
	"os"
)

var charset string
func init() {
        flag.StringVar(&charset, "c", "utf-8", "character set")
}

func main() {
	s, err := giawarc.CleanText(os.Stdin, charset)
	if err != nil {
		log.Fatal(err)
	}
	os.Stdout.WriteString(giawarc.CleanSpaces(s))
}
