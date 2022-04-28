package app

import (
	"github.com/gin-gonic/gin"
	"github.com/proactiongo/pagocore/di"
	"github.com/proactiongo/pagocore/ginsrv"
	log "github.com/sirupsen/logrus"
)

const logTag = "[pagocore.App] "

// NewApp creates new App instance.
// If ctn is nil, GetDefaultDIBuilder() will be called.
func NewApp(ctn *di.Container) *App {
	if ctn == nil {
		ctn = BuildDefaultContainer()
	}
	return &App{
		ctn: ctn,
	}
}

// PrepareRouterFn is a function to prepare router before run
type PrepareRouterFn func(router *gin.Engine, ctn *di.Container) error

// PrepareContainerFn is a function to prepare DI container before run
type PrepareContainerFn func(ctn *di.Container) error

// App is a main service app
type App struct {
	ctn *di.Container

	prepareCtnFn    PrepareContainerFn
	prepareRouterFn PrepareRouterFn

	initialized bool
}

// Init initializes App without starting the server
func (a *App) Init() {
	if a.initialized {
		return
	}
	a.initialized = true

	if a.prepareCtnFn != nil {
		err := a.prepareCtnFn(a.C())
		if err != nil {
			log.Fatal(logTag, "failed to prepare DI container: ", err)
		}
	}

	router := DIGetRouter(a.C())
	router.Use(ginsrv.M().SetDIContainer(a.C()))

	if a.prepareRouterFn != nil {
		err := a.prepareRouterFn(router, a.C())
		if err != nil {
			log.Fatal(logTag, "failed to init router: ", err)
		}
	}
}

// Run starts the App's server
func (a *App) Run() {
	a.Init()
	defer a.Close()

	router := DIGetRouter(a.ctn)
	conf := DIGetConfig(a.ctn)

	log.Info(logTag, "starting gin server")

	err := router.Run(":" + conf.Port)
	if err != nil {
		log.Fatal(err)
	}
}

// Close finalizes the App
func (a *App) Close() {
	a.ctn.Close()
}

// SetPrepareRouterFn sets init router hook
func (a *App) SetPrepareRouterFn(fn PrepareRouterFn) {
	a.prepareRouterFn = fn
}

// SetPrepareContainerFn sets prepare DI container hook
func (a *App) SetPrepareContainerFn(fn PrepareContainerFn) {
	a.prepareCtnFn = fn
}

// C returns an App's Container instance
func (a *App) C() *di.Container {
	if a.ctn == nil {
		log.Warn(logTag, "nil Container given to the App, switching to default")
		a.ctn = BuildDefaultContainer()
	}
	return a.ctn
}
