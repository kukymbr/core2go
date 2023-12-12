package di

import "errors"

// Builder errors
var (
	ErrDefinitionExists     = errors.New("definition already registered")
	ErrBuildFunctionMissing = errors.New("definition build function is missing")
	ErrDefinitionNotFound   = errors.New("definition not found")
)
