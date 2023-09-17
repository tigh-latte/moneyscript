package object

type Environment struct {
	*Environment
	s map[string]Object
}

func NewEnvironment(env *Environment) *Environment {
	return &Environment{
		Environment: env,
		s:           make(map[string]Object),
	}
}

func (e *Environment) Get(name string) (Object, bool) {
	o, ok := e.s[name]
	if !ok && e.Environment != nil {
		o, ok = e.Environment.Get(name)
	}
	return o, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.s[name] = val
	return val
}
