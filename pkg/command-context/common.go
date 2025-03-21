package commandcontext

// Command Context
type CC interface {
	Run(name string, args ...string) error
	CommandExists(name string) (bool, error)
}
