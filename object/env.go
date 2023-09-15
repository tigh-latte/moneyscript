package object

type Environment struct {
	s map[string]Object
}

func NewEnvironment() *Environment {
	return &Environment{
		s: make(map[string]Object),
	}
}

func (e *Environment) Get(name string) (Object, bool) {
	o, ok := e.s[name]
	return o, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.s[name] = val
	return val
}
