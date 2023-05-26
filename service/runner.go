package service

// Runner is an interface for a service's runner,
// such as gin or something else.
type Runner interface {
	// Run initializes the runner and runs it.
	Run() error
}
