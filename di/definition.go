package di

import (
	"fmt"
)

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

	obj   any
	built bool
}

// build builds dependency's object
func (d *Def) build(ctn *Container) (err error) {
	if d.built {
		return nil
	}

	d.built = true

	if d.Build == nil {
		return fmt.Errorf("%s: %w", d.Name, ErrBuildFunctionMissing)
	}

	d.obj, err = d.Build(ctn)

	if err != nil {
		d.obj = nil
	}

	return err
}

// definitions is a dependencies definitions map
type definitions map[string]Def

// ValidateFn is a dependency validation function
type ValidateFn func(ctn *Container) (err error)

// BuildFn is a dependency build function
type BuildFn func(ctn *Container) (obj any, err error)

// CloseFn is a dependency close function
type CloseFn func(obj any) (err error)
