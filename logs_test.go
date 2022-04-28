package pagocore_test

import (
	"github.com/proactiongo/pagocore"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitLogs(t *testing.T) {
	pagocore.Opt.LogLevelDft = log.WarnLevel
	pagocore.InitLogs()
	assert.Equal(t, log.WarnLevel, log.GetLevel())
	log.SetLevel(log.DebugLevel)
}
