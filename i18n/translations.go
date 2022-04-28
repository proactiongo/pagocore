package i18n

import (
	"bytes"
	"github.com/proactiongo/pagocore/utils"
	log "github.com/sirupsen/logrus"
	"text/template"
)

// Translations is a list of translations map in format:
// {KeyOrDftLangText: {Lang: Text}}
type Translations map[string]map[Language]Translation

// GetText returns text in specified lang
func (t Translations) GetText(textOrKey string, lang Language, tplData interface{}) string {
	tr := t.getTranslation(textOrKey, lang)
	return tr.GetText(tplData)
}

// getTranslation returns Translation object for the specified textOrKey in specified lang
func (t Translations) getTranslation(textOrKey string, lang Language) Translation {
	if lang == "" {
		lang = Source.DefaultLang
	}
	dft := Translation{
		Text: textOrKey,
	}

	texts, ok := t[textOrKey]
	if !ok {
		return dft
	}

	text, ok := texts[lang]
	if ok {
		return text
	}
	if lang != Source.DefaultLang {
		text, ok := texts[Source.DefaultLang]
		if ok {
			return text
		}
	}

	return dft
}

// Translation is a one translation item
type Translation struct {
	Text    string `json:"text" yaml:"text"`
	textTpl *template.Template
}

// GetText returns text with tplData applied
func (t *Translation) GetText(tplData interface{}) string {
	return t.getText(t.textTpl, t.Text, tplData)
}

// getText applies vars to the text template
func (t *Translation) getText(tpl *template.Template, text string, tplData interface{}) string {
	var err error

	if tpl == nil {
		tpl = template.New(utils.GenerateUUID())
		tpl, err = tpl.Parse(text)
		if err != nil {
			log.WithField("tpl_text", text).Warn("failed to parse text template: ", err)
			return text
		}
	}

	var b bytes.Buffer
	err = tpl.Execute(&b, tplData)
	if err != nil {
		log.WithFields(
			log.Fields{
				"tpl_text": text,
				"tpl_data": tplData,
			},
		).Warn("failed to execute text template: ", err)
		return text
	}
	return b.String()
}
