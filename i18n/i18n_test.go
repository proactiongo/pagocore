package i18n

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestT(t *testing.T) {
	readTestSource()

	var tx string

	tx = T("test_key_1", "", nil)
	assert.Equal(t, "Test English text 1", tx)

	tx = T("test_key_1", LangRu, nil)
	assert.Equal(t, "Тестовый текст на русском 1", tx)

	tx = T("Test English text 2", LangEn, nil)
	assert.Equal(t, "Test English text 2", tx)

	tx = T("Test English text 2", "zh", nil)
	assert.Equal(t, "Test English text 2", tx)

	tx = T("Test English text 3", LangEn, nil)
	assert.Equal(t, "Test English text 3 (override)", tx)

	testStruct := struct {
		Test string
	}{
		Test: "TheTest",
	}

	testMap := map[string]string{
		"Test": "TheTest",
	}

	tx = T("test_key_2_vars", LangEn, testMap)
	assert.Equal(t, "Test English text 4, with var: TheTest", tx)

	tx = T("test_key_2_vars", LangRu, testStruct)
	assert.Equal(t, "Тестовый текст на русском 4 с переменной: TheTest", tx)

	tx = T("_unknown_key_", LangEn, nil)
	assert.Equal(t, "_unknown_key_", tx)

	tx = T("test_key_3", LangRu, nil)
	assert.Equal(t, "Test English text 5", tx)

	tx = T("invalid_tpl", LangEn, nil)
	assert.Equal(t, "Invalid template {{.}", tx)
}

func TestNewSourceFromFile(t *testing.T) {
	_, err := NewSourceFromFile("i18n_test.yml")
	assert.NoError(t, err)

	_, err = NewSourceFromFile("_unknown_file_")
	assert.Error(t, err)

	_, err = NewSourceFromFile("source.go")
	assert.Error(t, err)
}

func TestLanguage_Filter(t *testing.T) {
	langs := map[string]Language{
		" en  ": LangEn,
		"EN\t":  LangEn,
		" ru ":  LangRu,
		"Ru":    LangRu,
	}

	for inp, expected := range langs {
		assert.Equal(t, expected, ParseLanguage(inp))
	}
}

func TestLanguage_Validate(t *testing.T) {
	valid := []string{
		" en ",
		"EN ",
		" RU",
		"ru\n\n",
	}
	invalid := []string{
		"",
		"  ",
		" !ok ",
		" ~~WTF** ",
		"notacodeatall",
	}

	for _, v := range valid {
		l := Language(v)
		assert.NoError(t, l.Validate())
	}

	for _, v := range invalid {
		l := Language(v)
		assert.ErrorIs(t, l.Validate(), ErrInvalidLang)
	}
}

func readTestSource() {
	var err error
	Source, err = NewSourceFromFile("i18n_test.yml")

	if err != nil {
		log.Fatal("failed to read i18n_test.yml: ", err)
	}
}
