package giawarc

import (
	"io"
	"strings"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/unicode/norm"
	"github.com/microcosm-cc/bluemonday"
)

var broken_content_types map[string]string
var policy *bluemonday.Policy

func init() {
	policy = bluemonday.StrictPolicy()

	broken_content_types = map[string]string {
		"txt": "text/plain",
		"text": "text/plain",
		"text/plan": "text/plain",
	}
}

func CleanContentType(content_type string) (string, string) {
	// TODO also return charset here
	sp := strings.SplitN(strings.ToLower(content_type), ";", 2)

	content_type = strings.TrimSpace(sp[0])
	fixed_content_type, ok := broken_content_types[content_type]
	if ok {
		content_type = fixed_content_type
	}

	// we have tags, and we *assume* that we only have charset=tag
	var charset string
	if len(sp) == 2 {
		if strings.HasPrefix(sp[1], "charset=") {
			charset = sp[1][8:]
		}
	}
	return content_type, charset
}


// Find the right reader to recode into UTF-8
func Recode(body io.Reader, charset string) io.Reader {
	if charset != "" {
		return body
	}

	enc, err := ianaindex.MIME.Encoding(charset)
	if err == nil {
		dec := enc.NewDecoder()
		return dec.Reader(body)
	}

	enc, err = ianaindex.IANA.Encoding(charset)
	if err == nil {
		dec := enc.NewDecoder()
		return dec.Reader(body)
	}

	return body
}

// Clean and sanitize the text, making sure the result is normalised
// UTF-8 free of any HTML markup
func CleanText(reader io.Reader, charset string) io.Reader {

	encoded   := Recode(reader, charset)       // transform to UTF-8
	normed    := norm.NFKC.Reader(encoded)     // normalise UTF-8
	sanitized := policy.SanitizeReader(normed) // strip out any HTML crap

	return sanitized
}
