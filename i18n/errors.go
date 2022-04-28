package i18n

import (
	"github.com/proactiongo/pagocore"
	"net/http"
)

// ErrInvalidLang is a Language validation error
var ErrInvalidLang = pagocore.NewError(http.StatusBadRequest, "invalid language code given")
