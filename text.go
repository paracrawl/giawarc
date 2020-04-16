package giawarc

// Lookup table for MIME types that suggest text
var text_types map[string]bool

type TextRecord struct {
	Source string
	Date string
	RecordId string
	URI string
	ContentType string
	Lang string
	Text string
}

type TextWriter interface {
	WriteText(*TextRecord) (int, error)
	Close() error
}

func init() {
	text_types = map[string]bool {
		"text/plain": true,
		"text/html": true,
		"text/vnd.wap.wml": true,
		"application/xml": true,
		"application/atom+xml": true,
		"application/opensearchdescription+xml": true,
		"application/rss+xml": true,
		"application/xhtml+xml": true,
	}
}

func IsText(content_type string) bool {
	_, ok := text_types[content_type]
	return ok
}
