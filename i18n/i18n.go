package i18n

import (
	"regexp"
	"strings"
)

// Available languages
const (
	LangRu = Language("ru")
	LangEn = Language("en")
)

// Source is a current texts source
var Source = &TextsSource{
	DefaultLang:  LangEn,
	Translations: make(Translations),
}

var langCodeRegx *regexp.Regexp

// Language is a language code
type Language string

// Filter value
func (l *Language) Filter() {
	s := strings.TrimSpace(l.String())
	s = strings.ToLower(s)
	*l = Language(s)
}

// Validate if value is ok
func (l *Language) Validate() error {
	l.Filter()
	if *l == "" {
		return ErrInvalidLang
	}
	if langCodeRegx == nil {
		langCodeRegx = regexp.MustCompile(`^([a-z]){2}([_-][a-z]{2,4})*$`)
	}
	if !langCodeRegx.MatchString(l.String()) {
		return ErrInvalidLang
	}
	return nil
}

// String returns Language as string
func (l *Language) String() string {
	return string(*l)
}

// T is a Translations.GetText alias
func T(textOrKey string, lang Language, tplData interface{}) string {
	return Source.T(textOrKey, lang, tplData)
}
