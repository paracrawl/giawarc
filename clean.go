package giawarc

import (
	"bytes"
	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/unicode/norm"
	"html"
	"io"
	"regexp"
	"strings"
)

var broken_content_types = map[string]string{
	"txt": "text/plain",
	"text": "text/plain",
	"text/plan": "text/plain",
}
var policy *bluemonday.Policy

var ms_re, ds_re *regexp.Regexp

func init() {
	policy = bluemonday.StrictPolicy()

	ms_re = regexp.MustCompile(`[ \t\r\p{Zs}\x{c2a0}]+`)
	ds_re = regexp.MustCompile(`[\p{Zs}]*[\x0a\x0b\x0c\x0d\p{Zl}\p{Zp}]+[\p{Zs}]*`)
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

func CleanSpaces(s string) string {
	s = strings.TrimSpace(s)
	s = ms_re.ReplaceAllLiteralString(s, " ")
	var buf strings.Builder

	for _, s := range(strings.Split(s, "\n")) {
		s := strings.TrimSpace(s)
		if len(s) == 0 {
			continue
		}
		buf.WriteString(s)
		buf.WriteString("\n")
	}
	//	s = ds_re.ReplaceAllLiteralString(s, "\n")
	//	return s
	return buf.String()
}

func FixInvalidUtf8(reader io.Reader) (string, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return "", err
	}

	var str strings.Builder
	// range has a neat side-effect of cleaning invalid utf-8
	for _, r := range buf.String() {
		_, err = str.WriteRune(r)
		if err != nil {
			return "", err
		}
	}

	return str.String(), nil
}

// Clean and sanitize the text, making sure the result is normalised
// UTF-8 free of any HTML markup
func CleanText(reader io.Reader, charset string) (string, error) {

	sanitized  := policy.SanitizeReader(reader) // strip out any HTML crap
	encoded    := Recode(sanitized, charset)    // transform to UTF-8
	valid, err := FixInvalidUtf8(encoded)     	// make sure all UTF-8 is valid
	if err != nil {
		return "", err
	}
	normed     := norm.NFKC.String(valid)       // normalise UTF-8
	unescaped  := html.UnescapeString(normed)   // take care of html &xx;

	return unescaped, nil
}
