package di

import (
	"errors"
	log "github.com/sirupsen/logrus"
)

// Container is a dependency container
type Container struct {
	defs DefsMap
}

// Has checks if dependency is registered in Container
func (c *Container) Has(name string) bool {
	_, ok := c.defs[name]
	return ok
}

// Get returns built dependency. Panics on error.
func (c *Container) Get(name string) interface{} {
	obj, err := c.SafeGet(name)
	if err != nil {
		panic(err)
	}
	return obj
}

// SafeGet returns built dependency
func (c *Container) SafeGet(name string) (interface{}, error) {
	def, ok := c.defs[name]
	if !ok {
		return nil, errors.New("[pagocore.di] dependency is not registered: " + name)
	}
	if def.Lazy {
		err := def.build(c)
		if err != nil {
			return nil, err
		}
		c.defs[name] = def
	}
	return def.obj, nil
}

// Close finalizes dependencies
func (c *Container) Close() {
	for _, def := range c.defs {
		if !def.built {
			continue
		}
		if def.Close == nil {
			continue
		}
		err := def.Close(def.obj)
		if err != nil {
			log.Error("[pagocore.di] failed to close dependency: ", err)
		}
	}
}
