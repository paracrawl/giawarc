package giawarc

// Lookup table for MIME types that suggest text
var text_types map[string]bool

type TextRecord struct {
	RecordId string
	URI string
	ContentType string
	Lang string
	Text string
}

type TextWriter interface {
	WriteText(*TextRecord) error
	Close() error
}

func init() {
	text_types = map[string]bool {
		"text/plain": true,
		"text/html": true,
	}
}

func IsText(content_type string) bool {
	_, ok := text_types[content_type]
	return ok
}
