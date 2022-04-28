package ginsrv

import (
	"github.com/gin-gonic/gin"
	"github.com/proactiongo/pagocore"
	"github.com/proactiongo/pagocore/di"
	"github.com/proactiongo/pagocore/i18n"
	"github.com/proactiongo/pagocore/tokens"
	"github.com/proactiongo/pagocore/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// NewContextHandler creates new ContextHandler instance
func NewContextHandler(ctx *gin.Context) *ContextHandler {
	return &ContextHandler{
		Context: ctx,
	}
}

// ContextHandler is a gin context wrapper
type ContextHandler struct {
	*gin.Context
}

// C returns gin.Context instance
func (h *ContextHandler) C() *gin.Context {
	return h.Context
}

// GetContainer returns di.Container instance
func (h *ContextHandler) GetContainer() *di.Container {
	ctn, ok := h.Get(KeyDIContainer)
	if !ok {
		panic("attempt to access non-initialized DI Container")
	}
	return ctn.(*di.Container)
}

// GetI18nSource returns i18n texts source from the container
func (h *ContextHandler) GetI18nSource() *i18n.TextsSource {
	return h.GetContainer().Get("pa_i18n").(*i18n.TextsSource)
}

// GetAccessClaims returns AccessTokenClaims from the current gin context
func (h *ContextHandler) GetAccessClaims() (*tokens.AccessTokenClaims, error) {
	c, ok := h.Get(KeyAccessClaims)
	if !ok {
		log.Warn("no claims initialized")
		return nil, pagocore.ErrTokenInvalid
	}
	claims, ok := c.(*tokens.AccessTokenClaims)
	if !ok {
		log.Warn("unexpected claims type")
		return nil, pagocore.ErrTokenInvalid
	}
	return claims, nil
}

// GetI18nLang gets language from gin context
func (h *ContextHandler) GetI18nLang() i18n.Language {
	v, ok := h.Get(KeyI18nLang)
	if !ok {
		return i18n.GetDefaultLang()
	}
	return i18n.ParseLanguage(v)
}

// SetI18nLang set lang for gix context
func (h *ContextHandler) SetI18nLang(lang i18n.Language) {
	h.Set(KeyI18nLang, lang)
}

// ExtractBearerToken extracts token from the 'Authorization: Bearer <token>' header
func (h *ContextHandler) ExtractBearerToken() (string, error) {
	return utils.ExtractBearerToken(h.Request)
}

// Err sends error as response
func (h *ContextHandler) Err(err error) {
	status := http.StatusInternalServerError
	e, ok := err.(*pagocore.Error)
	if ok {
		if e.Code >= 400 && e.Code <= 599 {
			status = e.Code
		}
	}
	h.ErrWithStatus(err, status)
}

// ErrS sends error message as response
func (h *ContextHandler) ErrS(msg string, status int) {
	h.Err(pagocore.NewError(status, msg))
}

// ErrWithStatus sends error as response with custom status
func (h *ContextHandler) ErrWithStatus(err error, status int) {
	e, ok := err.(*pagocore.Error)
	if !ok {
		e = pagocore.NewError(status, err.Error())
	}
	e.Localized = h.GetI18nSource().T(e.Error(), h.GetI18nLang(), nil)
	h.JSON(
		status,
		e,
	)
}
