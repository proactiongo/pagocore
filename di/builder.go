package di

import (
	"errors"
)

// Builder is a Container builder
type Builder struct {
	ctn *Container
	ord []string
}

// Add adds Def to the Container
func (b *Builder) Add(defs ...Def) error {
	b.initContainer()
	for _, def := range defs {
		if _, exists := b.ctn.defs[def.Name]; exists {
			return errors.New("[pagocore.di] definition with name `" + def.Name + "` already exists")
		}
		if def.Validate != nil {
			if err := def.Validate(b.ctn); err != nil {
				return err
			}
		}
		b.ctn.defs[def.Name] = def
		b.ord = append(b.ord, def.Name)
	}
	return nil
}

// Build prepares Container and builds non-lazy definitions
func (b *Builder) Build() (*Container, error) {
	b.initContainer()
	for _, name := range b.ord {
		def := b.ctn.defs[name]
		if !def.Lazy {
			if err := def.build(b.ctn); err != nil {
				return nil, err
			}
			b.ctn.defs[name] = def
		}
	}
	return b.ctn, nil
}

// initContainer creates container instance
func (b *Builder) initContainer() {
	if b.ctn == nil {
		b.ctn = &Container{
			defs: make(DefsMap),
		}
	}
}
