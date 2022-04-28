package i18n

import (
	"fmt"
	"strings"
)

// GetDefaultLang returns default Language
func GetDefaultLang() Language {
	return Source.DefaultLang
}

// ParseLanguage parses Language from inp (Language, string, []byte or fmt.Stringer)
func ParseLanguage(inp interface{}) Language {
	var lang Language
	var ok bool
	var err error

	lang, ok = inp.(Language)
	if ok {
		lang.Filter()
		if err = lang.Validate(); err == nil {
			return lang
		}
	}

	var b []byte
	var s fmt.Stringer
	var str string

	str, ok = inp.(string)
	if !ok {
		b, ok = inp.([]byte)
		if ok {
			str = string(b)
		} else {
			s, ok = inp.(fmt.Stringer)
			if ok {
				str = s.String()
			}
		}
	}

	if ok {
		lang = Language(str)
		lang.Filter()
	} else {
		lang = ""
	}

	if err = lang.Validate(); err != nil {
		return GetDefaultLang()
	}

	return lang
}

// ParseAcceptLanguageString does simple parsing of the first language
func ParseAcceptLanguageString(accept string) Language {
	parts := strings.Split(accept, ", ")
	p := strings.TrimSpace(parts[0])
	if p == "" {
		return GetDefaultLang()
	}

	parts = strings.Split(accept, "-")
	return ParseLanguage(parts[0])
}
