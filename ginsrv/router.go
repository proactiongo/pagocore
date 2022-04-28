package ginsrv

import (
	"github.com/gin-gonic/gin"
	"github.com/proactiongo/pagocore"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// GetDefaultRouter creates new gin router with default middlewares
func GetDefaultRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Recover panics with formatted log
	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		ctx := NewContextHandler(c)
		log.Error(recovered)
		if s, ok := recovered.(string); ok {
			ctx.ErrS(s, http.StatusInternalServerError)
		}
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}))

	// JSON-formatted logs
	router.Use(gin.LoggerWithFormatter(M().LogFormatter))
	router.Use(M().LogBody())

	// Init language from header
	router.Use(M().InitI18nLang())

	// 8 Mb limit for uploads
	router.MaxMultipartMemory = 8 << 20

	// Custom no-route error
	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(
			http.StatusNotFound,
			pagocore.NewError(http.StatusNotFound, "unknown endpoint"),
		)
	})

	return router
}
