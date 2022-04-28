package i18n

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

// NewSourceFromFile reads yaml file and creates new TextsSource instance
func NewSourceFromFile(path string) (*TextsSource, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Error(err)
		}
	}()
	b, err := ioutil.ReadAll(file)

	var source *TextsSource
	err = yaml.Unmarshal(b, &source)
	if err != nil {
		return nil, err
	}

	return source, nil
}

// TextsSource is an i18n source data
type TextsSource struct {
	DefaultLang  Language     `json:"default_lang" yaml:"default_lang"`
	Translations Translations `json:"translations" yaml:"translations"`
}

// T is a Translations.GetText alias
func (s *TextsSource) T(textOrKey string, lang Language, tplData interface{}) string {
	return s.Translations.GetText(textOrKey, lang, tplData)
}
