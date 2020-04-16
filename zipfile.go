package giawarc

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
)

var zip_types map[string]*regexp.Regexp

func init() {
	zip_types = map[string]*regexp.Regexp {
		"application/vnd.oasis.opendocument.text": regexp.MustCompile(`^content\.xml$`),
		"application/vnd.oasis.opendocument.spreadsheet": regexp.MustCompile(`^content\.xml$`),
		"application/vnd.oasis.opendocument.presentation": regexp.MustCompile(`^content\.xml$`),
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": regexp.MustCompile(`^word/document\.xml$`),
		"application/vnd.openxmlformats-officedocument.presentationml.presentation": regexp.MustCompile(`^ppt/slides/slide.*$`),
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": regexp.MustCompile(`^xl/sharedStrings\.xml$`),
		"application/epub+zip": regexp.MustCompile(`^.*ml$`),
	}
}

func IsZip(content_type, uri string) (string, bool) {
	if strings.HasSuffix(uri, "odt") {
		return "application/vnd.oasis.opendocument.text", true
	}
	if strings.HasSuffix(uri, "ods") {
		return "application/vnd.oasis.opendocument.spreadsheet", true
	}
	if strings.HasSuffix(uri, "odp") {
		return "application/vnd.oasis.opendocument.presentation", true
	}
	if strings.HasSuffix(uri, "docx") {
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document", true
	}
	if strings.HasSuffix(uri, "pptx") {
		return "application/vnd.openxmlformats-officedocument.presentationml.presentation", true
	}
	if strings.HasSuffix(uri, "xslx") {
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", true
	}
	if strings.HasSuffix(uri, "epub") {
		return "application/epub+zip", true
	}

	_, ok := zip_types[content_type]
	if ok {
		return content_type, true
	}

	return content_type, false
}

func ReadZipPayload(content_type string, body io.Reader) (buf io.Reader, err error){
	zipdata, err := ioutil.ReadAll(body)
	if err != nil {
		return
	}

	zip, err := zip.NewReader(bytes.NewReader(zipdata), int64(len(zipdata)))
	if err != nil {
		return
	}
	
	data := make([]byte, 0, 16384)

	fre := zip_types[content_type]
	for _, f := range zip.File {
		if ! fre.MatchString(f.Name) {
			continue
		}

		fp, err := f.Open()
		if err != nil {
			continue
		}

		contents, err := ioutil.ReadAll(fp)
		if err != nil {
			fp.Close()
			continue
		}

		fp.Close()

		data = append(data, contents...)
	}

	buf = bytes.NewBuffer(data)
	return
}
