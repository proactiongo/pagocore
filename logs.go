package pagocore

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

// Log custom fields
const (
	LogFieldTimestamp         = "@timestamp"
	LogFieldLevel             = "level"
	LogFieldMessage           = "message"
	LogFieldFunc              = "func"
	LogFieldService           = "service"
	LogFieldHostname          = "hostname"
	LogFieldAPIVersion        = "api_version"
	LogFieldType              = "log_type"
	LogFieldPath              = "path"
	LogFieldStatus            = "status"
	LogFieldMethod            = "method"
	LogFieldClientAppVersion  = "client_app_ver"
	LogFieldClientAppPlatform = "client_app_platform"
)

// Log types
const (
	LogTypeHTTPSrv = "http_srv"
	LogTypeHTTPIO  = "http_io"
	LogTypeApp     = "app"
)

// region WRITERS

// InitLogs initializes logs writers
func InitLogs() {
	log.SetLevel(Opt.LogLevelDft)
	log.SetFormatter(&log.JSONFormatter{
		FieldMap: log.FieldMap{
			log.FieldKeyTime:  LogFieldTimestamp,
			log.FieldKeyLevel: LogFieldLevel,
			log.FieldKeyMsg:   LogFieldMessage,
			log.FieldKeyFunc:  LogFieldFunc,
		},
	})

	log.SetOutput(initLogWriter(Opt.LogsPath + "/" + Opt.LogFileApp))
	log.AddHook(&LogDftFieldsHook{})

	gin.DefaultWriter = initLogWriter(Opt.LogsPath + "/" + Opt.LogFileGin)
}

// initLogWriter opens file log writer
func initLogWriter(path string) io.Writer {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Error("failed to open log file ", path)
		return os.Stdout
	}
	return io.MultiWriter(os.Stdout, f)
}

// endregion WRITERS

// region GIN BODY

// LogDftFieldsHook is a log's hook to add custom fields to all messages
type LogDftFieldsHook struct{}

// Levels to apply hook to
func (h *LogDftFieldsHook) Levels() []log.Level {
	return log.AllLevels
}

// Fire the hook
func (h *LogDftFieldsHook) Fire(e *log.Entry) error {
	_, ok := e.Data[LogFieldType]
	if !ok {
		e.Data[LogFieldType] = LogTypeApp
	}
	e.Data[LogFieldService] = Opt.ServiceName
	e.Data[LogFieldHostname] = Opt.GetHostname()
	e.Data[LogFieldAPIVersion] = Opt.APIVersion
	return nil
}

// endregion GIN BODY
