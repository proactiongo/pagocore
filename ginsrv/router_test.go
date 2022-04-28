package ginsrv

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDefaultRouter(t *testing.T) {
	_ = GetDefaultRouter()
	assert.Equal(t, gin.ReleaseMode, gin.Mode())
}
