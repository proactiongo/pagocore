package i18n_test

import (
	"github.com/proactiongo/pagocore"
	"github.com/proactiongo/pagocore/i18n"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseAcceptLanguageString(t *testing.T) {
	values := map[string]i18n.Language{
		"ru-RU, ru;q=0.9, en-US;q=0.8, en;q=0.7, fr;q=0.6": i18n.LangRu,
		"*":     i18n.LangEn,
		"en-US": i18n.LangEn,
		"":      i18n.LangEn,
	}

	for inp, expected := range values {
		lang := i18n.ParseAcceptLanguageString(inp)
		assert.Equal(t, expected, lang)
	}
}

func TestParseLanguage(t *testing.T) {
	inputs := []interface{}{
		"ru",
		[]byte("ru"),
		i18n.Language(" ru "),
		"",
		"WTF",
		&pagocore.Error{},
	}
	outputs := []i18n.Language{
		i18n.LangRu,
		i18n.LangRu,
		i18n.LangRu,
		i18n.LangEn,
		i18n.LangEn,
		i18n.LangEn,
	}

	for i, inp := range inputs {
		expected := outputs[i]
		lang := i18n.ParseLanguage(inp)
		assert.Equal(t, expected, lang)
	}
}
