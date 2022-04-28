package ginsrv

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/proactiongo/pagocore"
	"github.com/proactiongo/pagocore/di"
	"github.com/proactiongo/pagocore/i18n"
	"github.com/proactiongo/pagocore/tokens"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	// KeyAccessClaims is a param name for current access claims
	KeyAccessClaims = "PAAccessClaims"

	// KeyI18nLang is a context key for language
	KeyI18nLang = "PAI18nLang"

	// KeyDIContainer is a context key for a di.Container
	KeyDIContainer = "PADIContainer"
)

// M returns Middlewares instance
func M() *Middlewares {
	if middlewares == nil {
		middlewares = &Middlewares{}
	}
	return middlewares
}

// middlewares is a current Middlewares instance
var middlewares *Middlewares

// Middlewares contains middlewares functions
type Middlewares struct {
}

// SetDIContainer is a middleware to set app.App instance to the context
func (m *Middlewares) SetDIContainer(ctn *di.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(KeyDIContainer, ctn)
	}
}

// JWTAccess is an authorization by the Access token.
// Sets parsed claims to KeyAccessClaims param.
func (m *Middlewares) JWTAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := NewContextHandler(c)

		var err error
		sign, err := ctx.ExtractBearerToken()
		if err != nil {
			e := pagocore.NewError(http.StatusUnauthorized, err.Error())
			ctx.Err(e)
			ctx.Abort()
			return
		}

		err = m.initJWTAccess(ctx, sign)
		if err != nil {
			ctx.Err(err)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// NonRequiredJWTAccess is an authorization by the Access token if it is provided.
// Sets parsed claims to KeyAccessClaims param.
func (m *Middlewares) NonRequiredJWTAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := NewContextHandler(c)

		var err error
		sign, err := ctx.ExtractBearerToken()
		if err != nil {
			ctx.Next()
			return
		}

		err = m.initJWTAccess(ctx, sign)
		if err != nil {
			ctx.Err(err)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func (m *Middlewares) initJWTAccess(ctx *ContextHandler, sign string) error {
	claims, err := tokens.ParseAccessToken(sign)
	if err != nil {
		return err
	}
	ctx.Set(KeyAccessClaims, claims)

	err = claims.Valid()
	if err != nil {
		return err
	}

	if claims.Language != "" {
		ctx.SetI18nLang(claims.Language)
	}

	return nil
}

// RequireRole validates user role is equal to specified
func (m *Middlewares) RequireRole(role int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := NewContextHandler(c)
		claims, err := ctx.GetAccessClaims()
		if err != nil {
			ctx.Err(err)
			ctx.Abort()
			return
		}
		if claims.Role != role {
			ctx.Err(pagocore.ErrRoleNotAllowed)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// RequireRoleIn validates if user role is one of specified
func (m *Middlewares) RequireRoleIn(roles ...int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := NewContextHandler(c)
		claims, err := ctx.GetAccessClaims()
		if err != nil {
			ctx.Err(err)
			ctx.Abort()
			return
		}
		for _, role := range roles {
			if claims.Role == role {
				ctx.Next()
				return
			}
		}
		ctx.Err(pagocore.ErrRoleNotAllowed)
		ctx.Abort()
	}
}

// RequireRoleHigher validates if user role is equal or higher of specified
func (m *Middlewares) RequireRoleHigher(role int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := NewContextHandler(c)
		claims, err := ctx.GetAccessClaims()
		if err != nil {
			ctx.Err(err)
			ctx.Abort()
			return
		}
		if claims.Role < role {
			ctx.Err(pagocore.ErrRoleNotAllowed)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// RequireStrictServiceAllowList validates if access claims has a list of allowed services
func (m *Middlewares) RequireStrictServiceAllowList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := NewContextHandler(c)
		claims, err := ctx.GetAccessClaims()
		if err != nil {
			ctx.Err(err)
			ctx.Abort()
			return
		}

		if len(claims.GetAllowedServices()) == 0 {
			ctx.Err(pagocore.ErrTokenInvalid)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// InitI18nLang reads language from Accept-Language header and sets it to the context
func (m *Middlewares) InitI18nLang() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := NewContextHandler(c)
		accept := ctx.GetHeader("Accept-Language")
		lang := i18n.ParseAcceptLanguageString(accept)
		ctx.SetI18nLang(lang)
	}
}

// LogBody is a middleware to write response body to the log
func (m *Middlewares) LogBody() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		writer := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = writer
		ctx.Next()
		logger := log.WithFields(log.Fields{
			pagocore.LogFieldType:   pagocore.LogTypeHTTPIO,
			pagocore.LogFieldStatus: ctx.Writer.Status(),
			pagocore.LogFieldPath:   ctx.FullPath(),
			pagocore.LogFieldMethod: ctx.Request.Method,
		})
		if ctx.Writer.Status() >= 500 {
			logger.Error(writer.body.String())
		} else if ctx.Writer.Status() >= 400 {
			logger.Warn(writer.body.String())
		}
	}
}

// LogFormatter formats gin log record as JSON string
func (m *Middlewares) LogFormatter(param gin.LogFormatterParams) string {
	data := map[string]string{
		"@timestamp":            param.TimeStamp.Format(time.RFC3339),
		"ip":                    param.ClientIP,
		pagocore.LogFieldMethod: param.Method,
		pagocore.LogFieldPath:   param.Path,
		"proto":                 param.Request.Proto,
		pagocore.LogFieldStatus: strconv.FormatInt(int64(param.StatusCode), 10),
		"latency":               strconv.FormatFloat(param.Latency.Seconds(), 'f', 8, 64),
		"latency_fmt":           param.Latency.String(),
		"agent":                 param.Request.UserAgent(),
		"error":                 param.ErrorMessage,
		"request_body":          "-",
		"response_body_size":    strconv.FormatInt(int64(param.BodySize), 10),
	}

	if param.StatusCode >= 400 && param.Request.Body != nil {
		buf := new(strings.Builder)
		_, _ = io.Copy(buf, param.Request.Body)
		data["request_body"] = buf.String()
		defer func() {
			_ = param.Request.Body.Close()
		}()
	}

	data[pagocore.LogFieldType] = pagocore.LogTypeHTTPSrv
	data[pagocore.LogFieldService] = pagocore.Opt.ServiceName
	data[pagocore.LogFieldHostname] = pagocore.Opt.GetHostname()
	data[pagocore.LogFieldAPIVersion] = pagocore.Opt.APIVersion

	if param.StatusCode >= 500 {
		data[pagocore.LogFieldLevel] = "error"
	} else if param.StatusCode >= 400 {
		data[pagocore.LogFieldLevel] = "warning"
	} else {
		data[pagocore.LogFieldLevel] = "info"
	}

	data[pagocore.LogFieldClientAppVersion] = param.Request.Header.Get("client_app_version")
	data[pagocore.LogFieldClientAppPlatform] = param.Request.Header.Get("client_app_platform")

	// not using marshalling for speed-up
	values := make([]string, len(data))
	i := 0
	for key, val := range data {
		val = strings.ReplaceAll(val, `"`, `\"`)
		val = strings.TrimSpace(val)
		values[i] = `"` + key + `":"` + val + `"`
		i++
	}

	return "{" + strings.Join(values, ",") + "}\n"
}

// bodyLogWriter is a writer to write response body to the log
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write body
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
