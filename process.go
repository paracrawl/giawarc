package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"github.com/wolfgangmeyers/go-warc/warc"
	"github.com/wwaites/giawarc/cld2"
)

// This structure implements the reading side of the WARC preprocessor.
type WARCPreProcessor struct {
	wf *warc.WARCFile
	tw TextWriter

	TextRecords int  // records claiming to be text
	LangRecords int  // records where we can tell the language
	TotalRecords int // total records
	ContentCounts map[string]int // statistics about content types
	TextBytes  int // bytes claiming to be text
	LangBytes  int // bytes claiming to be text where we know the language
	TotalBytes int // total bytes
}

// Create a preprocessor given a readable buffer containing a (gzipped) WARC file.
// The second argument, the TextRecord channel is where texts that are found will
// be sent. It will be closed when the file is done.
func NewWARCPreProcessor(rc io.ReadCloser, tw TextWriter) (wp *WARCPreProcessor, err error) {
	var p WARCPreProcessor
	p.ContentCounts = make(map[string]int)
	p.wf, err = warc.NewWARCFile(rc)
	if err != nil {
		return
	}
	p.tw = tw
	wp = &p
	return
}

// Loop through each record and process it
func (p *WARCPreProcessor) Process() {
	reader := p.wf.GetReader()
	reader.Iterate(p.processRecord)
}

// Callback from the WARC reader Iterate function
func (p *WARCPreProcessor) processRecord(wr *warc.WARCRecord, err error) {
	if err != nil {
		// TODO log error here
		return
	}

	// content type of the WARC record not the payload
	content_type, _ := wr.GetHeader().Get("Content-Type")
	if content_type != "application/http; msgtype=response" {
		// nothing to do
		return
	}

	content_length := wr.GetHeader().GetContentLength()

	// record some statistics
	p.TotalRecords += 1
	p.TotalBytes   += content_length

	// get HTTP response out of the WARC file, and parse it
	payload := wr.GetPayload()
	resp, err := http.ReadResponse(bufio.NewReader(payload.GetReader()), nil)
	if err != nil {
		// TODO log error here
		return
	}

	content_type, charset := CleanContentType(resp.Header.Get("Content-Type"))

	// record some statistics
	count, ok := p.ContentCounts[content_type]
	if !ok {
		p.ContentCounts[content_type] = 1
	} else {
		p.ContentCounts[content_type] = count + 1
	}

	// here is where we would do, is it a PDF? transform to text and then continue,
	// is it a doc? transform to text and continue

	// If it is text...
	if IsText(content_type) {
		// record some statistics
		p.TextRecords += 1
		p.TextBytes   += content_length

		// transform to UTF-8 and normalise, strip HTML stuff
		body := CleanText(resp.Body, charset)

		// read the resulting document into RAM for language detection
		text, err := ioutil.ReadAll(body)
		if err != nil {
			// TODO log error
			return
		}

		//text = strings.TrimSpace(text)
		//text = strings.ReplaceAll(text, "\n", " ")
		s_text := string(text)
		lang, ok := cld2.DetectLang(s_text)
		if !ok {
			return
		}

		// record some statistics
		p.LangRecords += 1
		p.LangBytes   += content_length

		uri, _ := wr.GetHeader().Get("WARC-Target-URI")

		// send off a TextRecord to whatever will write it
		rec := TextRecord{
			RecordId: wr.GetHeader().GetRecordId(),
			URI: uri,
			ContentType: content_type,
			Lang: lang,
			Text: s_text,
		}

		p.tw.WriteText(&rec)
	}
}

func (p *WARCPreProcessor) ContentTypeStats() ContentStats {
	cts := make(ContentStats, len(p.ContentCounts))
	for k, v := range p.ContentCounts {
		s := ContentStat{ ContentType: k, Prevalence: float64(v)/float64(p.TotalRecords) }
		cts = append(cts, s)
	}
	sort.Sort(cts)
	return cts
}

type ContentStat struct {
	ContentType string
	Prevalence float64
}

type ContentStats []ContentStat

func (cts ContentStats) Len() int {
	return len(cts)
}

func (cts ContentStats) Less(i, j int) bool {
	return cts[i].Prevalence > cts[j].Prevalence
}

func (cts ContentStats) Swap(i, j int) {
	tmp := cts[i]
	cts[i] = cts[j]
	cts[j] = tmp
}
