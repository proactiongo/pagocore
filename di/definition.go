package di

import "errors"

// DefsMap is a dependencies definitions map
type DefsMap map[string]Def

// ValidateFn is a dependency validation function
type ValidateFn func(ctn *Container) error

// BuildFn is a dependency build function
type BuildFn func(ctn *Container) (interface{}, error)

// CloseFn is a dependency close function
type CloseFn func(obj interface{}) error

// Def is a dependency definition
type Def struct {
	// Name is a dependency name
	Name string

	// Lazy is a flag. If true, Build will be executed only on Container.Get() call.
	Lazy bool

	// Validate validates dependency definition on add
	Validate ValidateFn

	// Build builds dependency object
	Build BuildFn

	// Close finalizes dependency object
	Close CloseFn

	obj   interface{}
	built bool
}

// build builds dependency's object
func (d *Def) build(ctn *Container) error {
	if d.built {
		return nil
	}
	d.built = true

	if d.Build == nil {
		return errors.New("[pagocore.di] definition `" + d.Name + "`: build function is not defined")
	}

	var err error
	d.obj, err = d.Build(ctn)
	return err
}
