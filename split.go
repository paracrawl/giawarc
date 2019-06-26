package giawarc

import (
	"github.com/wwaites/giawarc/nbp"
	"strings"
	"regexp"
)

var cjk_lang = map[string]bool {
	"yue": true,
	"zh": true,
}

var end1_re, end2_re, end3_re *regexp.Regexp
var cjk1_re, cjk2_re *regexp.Regexp
var dword_re, acro_re, num_re *regexp.Regexp

func init() {
	// See https://github.com/kpu/preprocess/blob/master/moses/ems/support/split-sentences.perl

	// Non-period end of sentence markers (?!) followed by sentence starters.
	end1_re = regexp.MustCompile(`([?!]) +([¿¡'"\(\[\p{Pi}]* *[\p{Lu}])`)
	// ulti-dots or punctuation followed by sentence starters.
	end2_re = regexp.MustCompile(`(\.\.+|[?!\.] *[\'\"\)\]\p{Pf}]+) +([¿¡'"\(\[\p{Pi}]* *[\p{Lu}])`)
	// sentences that end with some sort of punctuation, and are followed by a
	// sentence starter punctuation and upper case.
	end3_re = regexp.MustCompile(`([?!\.]) +([¿¡'"\(\[\p{Pi}]+ *[\p{Lu}])`)

	// Chinese uses unusual end-of-sentence markers. These are NOT followed by
	// whitespace.  Nor is there any idea of capitalization.
	cjk1_re = regexp.MustCompile(`([。．？！♪])`)
	// Western end of sentence followed by ideograph is always end of sentence
	cjk2_re = regexp.MustCompile(`([\.?!]) *(\p{Han}|\p{Hangul}|\p{Hiragana}|\p{Katakana}|\p{Yi})`)

	dword_re = regexp.MustCompile(`([0-9A-Za-z\.\-]*)([\'\"\)\]\%\p{Pf}]*)\.+$`)
	acro_re  = regexp.MustCompile(`(\.)[\p{Lu}]+(\.+)$`)
	num_re   = regexp.MustCompile(`^[0-9]+[\.?!\'\"\)\]\p{Pf}]*`)
}

func SplitSentences(s, lang string) string {
	s = strings.TrimSpace(s)
	s = end1_re.ReplaceAllString(s, "${1}\n${2}")
	s = end2_re.ReplaceAllString(s, "${1}\n${2}")
	s = end3_re.ReplaceAllString(s, "${1}\n${2}")

	// special handling for Chinese, Japanese, Korean
	_, iscjk := cjk_lang[lang]
	if iscjk {
		s = cjk1_re.ReplaceAllString(s, "${1}\n${2}")
		s = cjk2_re.ReplaceAllString(s, "${1}\n${2}")
	}

	nb := nbp.NBP[lang]
	nu := nbp.NUM[lang]

	var buf strings.Builder
	words := strings.Split(s, " ")
	for i, word := range(words) {
		buf.WriteString(word)

		// does the string end with one or more dots?
		parts := dword_re.FindStringSubmatch(word)
		if parts == nil {
			buf.WriteString(" ")
			continue
		}

		// if a known honorific and not followed by extra punctuation do not break
		if nb != nil {
			_, ok := nb[parts[1]]
			if ok {
				buf.WriteString(" ")
				continue
			}
		}

		// if an acronym, do not break
		if acro_re.MatchString(word) {
			buf.WriteString(" ")
			continue
		}

		// if it is a special word followed by only numbers
		if nu != nil {
			_, ok := nu[parts[1]]
			if ok && i < len(words) - 1 && num_re.MatchString(words[i+1]) {
				buf.WriteString(" ")
				continue
			}
		}

		// none of the above, newline.
		buf.WriteString("\n")
	}
	return strings.TrimSpace(buf.String())
}
