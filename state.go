package gocket

import "fmt"

func (g *Gocket) ManageState(name string, value any) {
	g.state[name] = value
}

// State returns nil if the state is not present
func (ctx *GocketCtx) State(name string) any {
	return ctx.gocket.state[name]
}

func GetState[T any](ctx *GocketCtx, state string) T {
	s := ctx.State(state)
	if s == nil {
		panic(fmt.Sprintf("State was not present: %s", state))
	}
	val, ok := s.(T)
	if !ok {
		// This variable exists for printing T
		var t T
		panic(fmt.Sprintf("State \"%s\" was not of type: %T but of type: %T", state, t, s))
	}
	return val
}
