package pagocore_test

import (
	"github.com/proactiongo/pagocore"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetServiceInfo(t *testing.T) {
	info := pagocore.GetServiceInfo()
	assert.Equal(t, pagocore.Opt.APIVersion, info.Version)
	assert.Equal(t, pagocore.Opt.APIBasePath, info.BasePath)
}

func TestOptions_GetHostname(t *testing.T) {
	opt := &pagocore.Options{
		Hostname: "test_hostname",
	}
	assert.Equal(t, "test_hostname", opt.GetHostname())

	osHost, _ := os.Hostname()

	opt = &pagocore.Options{}
	assert.Equal(t, osHost, opt.GetHostname())
}

func TestReadConfig(t *testing.T) {
	pagocore.Opt.ConfigFilePath = ".test.env"
	pagocore.Opt.ConfigFileType = "env"

	conf, err := pagocore.ReadConfig()
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "test1", conf.GetString("test_env_var"))
	assert.Equal(t, "test2", conf.GetString("OTHER_TEST_ENV_VAR"))

	pagocore.Opt.ConfigFilePath = "__unkown_file__.env"
	_, err = pagocore.ReadConfig()
	assert.Error(t, err)
}
