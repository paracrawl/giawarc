// Package cld2 implements language detection using the
// Compact Language Detector.
//
// This package includes the relevant sources from the cld2
// project, so it doesn't require any external dependencies.
// For more information about CLD2, see https://code.google.com/p/cld2/.
package cld2

// #cgo CXXFLAGS: -std=c++98 -Wall -O2
// #cgo LDFLAGS: -lstdc++
// #include <stdlib.h>
// #include "cld2.h"
import "C"
import "unsafe"

// UnknownLanguage language code for unknown languages.
const UnknownLanguage = "un"

// UnknownLanguageName label for unknown languages.
const UnknownLanguageName = "UNKNOWN_LANGUAGE"

// DetectLang returns the language code for detected language
// in the given text.
func DetectLang(text string) (string, bool) {
	cs := C.CString(text)
	defer C.free(unsafe.Pointer(cs))

	reliable := C.int(0)

	res := C.DetectLang(cs, -1, &reliable)

	lang := UnknownLanguage
	if res != nil {
		lang = C.GoString(res)
	}
	return lang, reliable == 1
}

// LanguageNameFromCode returns a human readable language name
func LanguageNameFromCode(code string) string {
	if code == "" || code == UnknownLanguage {
		return UnknownLanguageName
	}
	cs := C.CString(code)
	res := C.LanguageNameFromCode(cs)
	C.free(unsafe.Pointer(cs))

	name := UnknownLanguageName
	if res != nil {
		name = C.GoString(res)
	}
	return name
}
