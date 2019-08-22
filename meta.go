package giawarc

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"github.com/paracrawl/go-warc/warc"
	"github.com/paracrawl/giawarc/cld2"
)

type WARCMetaProcessor struct {
	wf *warc.WARCFile
	Filename string
}

func NewWARCMetaProcessor(rc io.ReadCloser, filename string) (wp *WARCMetaProcessor, err error) {
	var p WARCMetaProcessor
	p.Filename = filename
	p.wf, err = warc.NewWARCFile(rc)
	if err != nil {
		return
	}
	wp = &p
	return
}

// Loop through each record and process it
func (p *WARCMetaProcessor) Process() {
	reader := p.wf.GetReader()
	reader.Iterate(p.processRecord)
}

func (p *WARCMetaProcessor) processRecord(wr *warc.WARCRecord, err error) {
	if err != nil {
		if err != io.EOF {
			log.Printf("Error reading WARC record: %v", err)
		}
		return
	}

//	warc_type := wr.GetHeader().GetType()
//	if warc_type == "warcinfo" {
//		p.Filename, _ = wr.GetHeader().Get("WARC-Filename")
//	}

	// content type of the WARC record not the payload
	content_type, _ := wr.GetHeader().Get("Content-Type")
	if !strings.Contains(content_type,"application/http") || !strings.Contains(content_type,"response") {
		return
	}

	uri, _ := wr.GetHeader().Get("WARC-Target-URI")

	content_length := wr.GetHeader().GetContentLength()
	date, _ := wr.GetHeader().Get("WARC-Date")

	// get HTTP response out of the WARC file, and parse it
	payload := wr.GetPayload()
	resp, err := http.ReadResponse(bufio.NewReader(payload.GetReader()), nil)
	reader := payload.GetReader()
	charset := "utf-8"
	if err != nil {
		return
	}

	content_type, charset = CleanContentType(resp.Header.Get("Content-Type"))

	var lang string
	if IsText(content_type) {
		text, err := CleanText(reader, charset)
		if err == nil {
			lang, _ = cld2.DetectLang(text)
		} else {
			log.Printf("Error reading content for %v: %v", uri, err)
		}
	}

	recid := wr.GetHeader().GetRecordId()

	fmt.Printf("%s\t%s\t%s\t%s\t%v\t%s\t%s\n",
		p.Filename, recid, uri, date, content_length, content_type, lang)
}

