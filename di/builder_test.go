package di_test

import (
	"errors"
	"github.com/proactiongo/pagocore/di"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuilder_Build(t *testing.T) {
	b := &di.Builder{}

	err := b.Add(
		di.Def{
			Name: "test1",
			Build: func(ctn *di.Container) (interface{}, error) {
				return "testval1", nil
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
	)
	assert.NoError(t, err)

	err = b.Add(
		di.Def{
			Name: "test1",
		},
	)
	assert.Error(t, err)

	err = b.Add(
		di.Def{
			Name: "test2",
			Validate: func(ctn *di.Container) error {
				return errors.New("expected error")
			},
		},
	)
	assert.Error(t, err)

	err = b.Add(
		di.Def{
			Name: "test3",
			Build: func(ctn *di.Container) (interface{}, error) {
				return nil, errors.New("expected error")
			},
			Lazy: true,
		},
	)
	assert.NoError(t, err)

	err = b.Add(
		di.Def{
			Name: "test4",
			Build: func(ctn *di.Container) (interface{}, error) {
				return "testval4", nil
			},
		},
	)
	assert.NoError(t, err)

	err = b.Add(
		di.Def{
			Name: "test5",
			Build: func(ctn *di.Container) (interface{}, error) {
				if !ctn.Has("test1") {
					return nil, errors.New("unexpected build order")
				}
				return "testval5", nil
			},
			Lazy: true,
			Close: func(obj interface{}) error {
				return errors.New("expected error")
			},
		},
	)
	assert.NoError(t, err)

	err = b.Add(
		di.Def{
			Name: "test6",
			Build: func(ctn *di.Container) (interface{}, error) {
				return "testval6", nil
			},
			Lazy: true,
		},
	)
	assert.NoError(t, err)

	ctn, err := b.Build()
	if !assert.NoError(t, err) {
		return
	}

	_, err = ctn.SafeGet("test3")
	assert.Error(t, err)

	v := ctn.Get("test1")
	assert.Equal(t, "testval1", v)

	v = ctn.Get("test5")
	assert.Equal(t, "testval5", v)

	assert.Panics(t, func() {
		ctn.Get("unknown_key")
	})

	ctn.Close()

	b = &di.Builder{}
	err = b.Add(
		di.Def{
			Name: "test7",
			Build: func(ctn *di.Container) (interface{}, error) {
				return nil, errors.New("expected error")
			},
		},
	)
	assert.NoError(t, err)
	_, err = b.Build()
	assert.Error(t, err)

	b = &di.Builder{}
	err = b.Add(
		di.Def{
			Name: "test8",
		},
	)
	assert.NoError(t, err)
	_, err = b.Build()
	assert.Error(t, err)
}
