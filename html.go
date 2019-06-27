package giawarc

import (
	"bytes"
	"errors"
	"golang.org/x/net/html"
	"io"
	"strings"
)

var startNL = map[string]bool {
	"ul": true,
	"ol": true,
	"dl": true,
	"tr": true,
}

var endNL = map[string]bool {
	"p": true,
	"div": true,
	"span": true,
	"li": true,
	"dd": true,
	"th": true,
	"td": true,
	"h1": true,
	"h2": true,
	"h3": true,
	"h4": true,
	"h5": true,
	"h6": true,
	"h7": true,
	"h8": true,
	"h9": true,
}

var selfNL = map[string]bool {
	"br": true,
}

var noText = map[string]bool {
	"script": true,
	"noscript": true,
	"style": true,
}

func HtmlToText(r io.Reader) (b *bytes.Buffer, err error) {
	var buf bytes.Buffer
	var lastTok string

	tokenizer := html.NewTokenizer(r)

	for {
		if tokenizer.Next() == html.ErrorToken {
			err = tokenizer.Err()
			if err == io.EOF {
				// End of input means end of processing
				return &buf, nil
			}
			// Raw tokenizer error
			return
		}

		token := tokenizer.Token()
		switch token.Type {
		case html.DoctypeToken:
		case html.CommentToken:
		case html.StartTagToken:
			if _, ok := startNL[token.Data]; ok {
				buf.WriteString("\n")
			}
			lastTok = token.Data
//			fmt.Fprintf(&buf, "<<<%s>>>\n", token.Data)
//			buf.WriteString(token.Data)
		case html.EndTagToken:
			if _, ok := endNL[token.Data]; ok {
				buf.WriteString("\n")
			} else {
				buf.WriteString(" ")
			}
		case html.SelfClosingTagToken:
			if _, ok := selfNL[token.Data]; ok {
				buf.WriteString("\n")
			}
		case html.TextToken:
			if _, ok := noText[lastTok]; !ok {
				buf.WriteString(strings.ReplaceAll(token.Data, "\n", " "))
			}
		default:
			// A token that didn't exist in the html package when we wrote this
			return nil, errors.New("unknown token")
		}
	}

	b = &buf
	return
}
