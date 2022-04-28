package i18n

// Translatable are structures which can translate their values using the global Source
type Translatable interface {
	// Translate struct values to the specified lang
	Translate(lang Language)
}

// TranslatableT are structures which can translate their values using the specified TextsSource
type TranslatableT interface {
	// Translate struct values to the specified lang
	Translate(texts *TextsSource, lang Language)
}
