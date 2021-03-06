package giawarc

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/paracrawl/giawarc/cld2"
//	"github.com/jmhodges/gocld3/cld3"
	"github.com/paracrawl/go-warc/warc"
	"github.com/spaolacci/murmur3"
)

// This structure implements the reading side of the WARC preprocessor.
type WARCPreProcessor struct {
	wf *warc.WARCFile
	tw TextWriter

	inputHashes map[uint32]struct{}
	outputHashes map[uint32]struct{}
	inputHashReader GzOrXzReader
	outputHashWriter ZipWriter
	inputHashing bool
	outputHashing bool

	langDetection string
//	cld3Model cld3.LanguageIdentifier

	Filename      string
	TextRecords   int            // records claiming to be text
	LangRecords   int            // records where we can tell the language
	TotalRecords  int            // total records
	ContentCounts map[string]int // statistics about content types
	TextBytes     int            // bytes claiming to be text
	LangBytes     int            // bytes claiming to be text where we know the language
	TotalBytes    int            // total bytes
}

// Create a preprocessor given a readable buffer containing a (gzipped) WARC file.
// The second argument, the TextRecord channel is where texts that are found will
// be sent. It will be closed when the file is done.
func NewWARCPreProcessor(rc io.ReadCloser, tw TextWriter, inputHashReader GzOrXzReader, outputHashWriter ZipWriter, langDetection string) (wp *WARCPreProcessor, err error) {
	var p WARCPreProcessor
	p.ContentCounts = make(map[string]int)
	p.wf, err = warc.NewWARCFile(rc)
	if err != nil {
		return
	}
	p.tw = tw

	p.inputHashReader = inputHashReader
	p.outputHashWriter = outputHashWriter
	p.inputHashing = (inputHashReader != (GzOrXzReader{}))
	p.outputHashing = (outputHashWriter != (ZipWriter{}))
	p.inputHashes = make(map[uint32]struct{})
	p.outputHashes = make (map[uint32]struct{})
	p.langDetection = langDetection
	wp = &p
	return
}

// Loop through each record and process it
func (p *WARCPreProcessor) Process() {
	if p.inputHashing {
		p.inputHashes = p.inputHashReader.ReadHashes()
	}
/*
	if p.langDetection == "cld3" {
		langIdModel, err := cld3.NewLanguageIdentifier(0,1024)
		if err != nil {
			log.Println("Error creating cld3 model")
			return
		}
		p.cld3Model = langIdModel
		defer cld3.FreeLanguageIdentifier(p.cld3Model)
  }
*/
	reader := p.wf.GetReader()
	reader.Iterate(p.processRecord)
	if p.outputHashing {
		p.outputHashWriter.WriteHashes(p.outputHashes)
	}
}

// Callback from the WARC reader Iterate function
func (p *WARCPreProcessor) processRecord(wr *warc.WARCRecord, err error) {
	if err != nil {
		if err != io.EOF {
			log.Printf("Error reading WARC record: %v", err)
		}
		return
	}

	warc_type := wr.GetHeader().GetType()
	if warc_type == "warcinfo" {
		p.Filename, _ = wr.GetHeader().Get("WARC-Filename")
		log.Printf("Processing %v", p.Filename)
	}

	// content type of the WARC record not the payload
	content_type, _ := wr.GetHeader().Get("Content-Type")
	if !strings.Contains(content_type, "application/http") || !strings.Contains(content_type, "response") {
		// nothing to do
		// log.Printf("Ignoring WARC record (not application/http)")
		return
	}

	content_length := wr.GetHeader().GetContentLength()

	date, _ := wr.GetHeader().Get("WARC-Date")

	// record some statistics
	p.TotalRecords += 1
	p.TotalBytes += content_length

	uri, _ := wr.GetHeader().Get("WARC-Target-URI")
	// skip robots.txt
	if strings.HasSuffix(uri, "robots.txt") {
		return
	}

	// get HTTP response out of the WARC file, and parse it
	payload := wr.GetPayload()
	resp, err := http.ReadResponse(bufio.NewReader(payload.GetReader()), nil)
	reader := payload.GetReader()

	if err == nil {
		content_type, _ = CleanContentType(resp.Header.Get("Content-Type"))

		// record some statistics
		count, ok := p.ContentCounts[content_type]
		if !ok {
			p.ContentCounts[content_type] = 1
		} else {
			p.ContentCounts[content_type] = count + 1
		}

		// here is where we would do, is it a PDF? transform to text and then continue,
		// is it a doc? transform to text and continue
		content_type, iszip := IsZip(content_type, uri)
		if iszip {
			reader, err = ReadZipPayload(content_type, resp.Body)
			if err != nil {
				return
			}
		} else if IsText(content_type) {
			reader = resp.Body
		} else {
			return
		}
	}

	p.TextRecords += 1
	p.TextBytes += content_length
	// transform to UTF-8 and normalise, strip HTML stuff
	text, err := CleanText(reader, content_type)
	if err != nil {
		log.Printf("Error reading content for %v: %v (read %d bytes)", uri, err, len(text))
		return
	}

	var lang string
	var ok bool
/*
	if p.langDetection == "cld3" {
		res := p.cld3Model.FindLanguage(text)
		lang = res.Language
		ok = res.IsReliable
	} else {
*/
		lang, ok = cld2.DetectLang(text)
/*
	}
*/
	if !ok {
		return
	}

	tidied := CleanSpaces(text) // clean up excess whitespace

	if p.outputHashing || p.inputHashing {
		// hash clean text
		newhash := murmur3.Sum32([]byte (tidied))
		_, exists := p.inputHashes[newhash]
		if exists { return }
		_, exists = p.outputHashes[newhash]
		if exists { return }
		// store new hash and continue
		p.outputHashes[newhash] = Empty
	}

	// record some statistics
	p.LangRecords += 1
	p.LangBytes += content_length

	recid := wr.GetHeader().GetRecordId()
	recid = strings.TrimPrefix(recid, "<urn:uuid:")
	recid = strings.TrimSuffix(recid, ">")

	// send off a TextRecord to whatever will write it
	rec := TextRecord{
		Source:      p.Filename,
		Date:        date,
		RecordId:    recid,
		URI:         uri,
		ContentType: content_type,
		Lang:        lang,
		Text:        tidied,
	}

	_, err = p.tw.WriteText(&rec)
}

// Utility to get statistics about content types for printing out.
func (p *WARCPreProcessor) ContentTypeStats() ContentStats {
	cts := make(ContentStats, len(p.ContentCounts))
	for k, v := range p.ContentCounts {
		s := ContentStat{ContentType: k, Prevalence: float64(v) / float64(p.TotalRecords)}
		cts = append(cts, s)
	}
	sort.Sort(cts)
	return cts
}

// A statistic about a content type
type ContentStat struct {
	ContentType string
	Prevalence  float64
}

// An array of content-type statistics, with a new name so that we
// can sort it by prevalence
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
