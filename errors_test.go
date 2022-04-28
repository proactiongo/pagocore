package pagocore_test

import (
	"github.com/proactiongo/pagocore"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNewError(t *testing.T) {
	err := pagocore.NewError(400, "test")
	assert.Error(t, err)
	assert.Equal(t, 400, err.Code)
	assert.Equal(t, "test", err.Message)

	err = pagocore.NewError(http.StatusNotFound)
	assert.Error(t, err)
	assert.Equal(t, 404, err.Code)
	assert.Equal(t, "Not Found", err.Message)
}

func TestError_Error(t *testing.T) {
	errs := map[string]*pagocore.Error{
		"test1": pagocore.NewError(1, "test1"),
		"test2": pagocore.NewError(2, "test2"),
		"test3": pagocore.NewError(2, "test", 3),
	}

	for text, err := range errs {
		assert.Equal(t, text, err.Error())
	}
}
