package pagocore

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

// Opt shares package options
var Opt = &Options{
	ConfigFilePath: "./.env",
	ConfigFileType: "env",

	APIVersion:  "v1",
	APIBasePath: "/api/v1",

	LogsPath:    "data/logs",
	LogFileGin:  "gin.log",
	LogFileApp:  "app.log",
	LogLevelDft: log.InfoLevel,
}

// Options represents package options
type Options struct {
	// ConfigFilePath is a path to an app config file
	ConfigFilePath string
	// ConfigFileType is an app config file type, e. g. "env"
	ConfigFileType string

	// LogsPath is a path to the log files directory
	LogsPath string
	// LogFileGin is a filename of gin server log
	LogFileGin string
	// LogFileApp is a filename of app server log
	LogFileApp string
	// LogLevelDft is a default log level if none is defined in config file
	LogLevelDft log.Level

	// APIVersion is an API version identifier, e. g. "v1"
	APIVersion string
	// APIBasePath is a service API base path, e. g. /api/v1
	APIBasePath string

	// ServiceName is a service identifier
	ServiceName string
	// ServiceDescription is a human-readable service description
	ServiceDescription string
	// ServiceRepo is an URL to service git repository
	ServiceRepo string

	// Hostname of current node if is required to override os.Hostname() value
	Hostname string

	// JWTPassword is JWT password key
	JWTPassword []byte
}

// GetHostname returns hostname from options or OS
func (o *Options) GetHostname() string {
	if o.Hostname == "" {
		o.Hostname, _ = os.Hostname()
	}
	return o.Hostname
}

// ServiceInfo represents info about the service
type ServiceInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	BasePath    string `json:"base_path"`
	Description string `json:"description"`
	Repo        string `json:"repo"`
}

// GetServiceInfo returns map with service info
func GetServiceInfo() *ServiceInfo {
	return &ServiceInfo{
		Name:        Opt.ServiceName,
		Version:     Opt.APIVersion,
		BasePath:    Opt.APIBasePath,
		Description: Opt.ServiceDescription,
		Repo:        Opt.ServiceRepo,
	}
}

// ReadConfig loads configuration from the config file to the viper instance
func ReadConfig() (*viper.Viper, error) {
	conf := viper.New()

	conf.SetConfigFile(Opt.ConfigFilePath)
	conf.SetConfigType(Opt.ConfigFileType)

	err := conf.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return conf, nil
}
