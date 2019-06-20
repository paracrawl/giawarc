package cld2

import "testing"

func TestLanguageNameFromCode_withUnknown_success(t *testing.T) {
	name := LanguageNameFromCode(UnknownLanguage)
	if UnknownLanguageName != name {
		t.Fatalf("Expected unknown language, found %s.", name)
	}
}

func TestLanguageNameFromCode_withEn_success(t *testing.T) {
	name := LanguageNameFromCode("en")
	if "ENGLISH" != name {
		t.Fatalf("Expected ENGLISH, found %s.", name)
	}
}

func TestDetectLang_withShortText_unreliable(t *testing.T) {
	lang, reliable := DetectLang("banana")
	if reliable {
		t.Fatalf("Expected results to be unreliable, found %v.", reliable)
	}
	if "un" != lang {
		t.Fatalf("Expected 'un', found %s.", lang)
	}
}

func TestDetectLang_withPersian_success(t *testing.T) {
	lang, reliable := DetectLang(`رُم پایتخت کشور ایتالیا، بزرگترین و پرجمعیت‌ترین شهر این کشور با ۲۶۴۹۷۲۴ سکن
	ه و همچنین مرکز ناحیهٔ لاتزیو می‌باشد. هم چنین رم با مساحت ۱۳۶۲۸۷ `)
	if !reliable {
		t.Fatalf("Expected results to be reliable, found %v.", reliable)
	}
	if "fa" != lang {
		t.Fatalf("Expected 'fa', found %s.", lang)
	}
}
