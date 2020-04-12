package object

// Environment the environment struct
type Environment struct {
	store map[string]Object
}

// NewEnvironment creates and returns a new environment
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

// Get returns value associated with the given environment key (name)
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

// Set associates the value with the given environment key (name) in the environment
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
