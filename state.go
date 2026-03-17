package gocket

import (
	"fmt"
	"sync"
)

type SafeState[T any] struct {
	mutex sync.Mutex
	val   T
}

func (s *SafeState[T]) Set(val T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.val = val
}

func (s *SafeState[T]) Get() T {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	val := s.val
	return val
}

func (s *SafeState[T]) Update(fn func(*T)) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	fn(&s.val)
}

func (g *Gocket) ManageState(name string, value any) {
	g.state[name] = value
}

// State returns nil if the state is not present
func (ctx *GocketCtx) State(name string) any {
	state, ok := ctx.localState[name]
	if ok {
		return state
	}
	return ctx.gocket.state[name]
}

func (ctx *GocketCtx) SetLocalState(name string, val any) {
	ctx.localState[name] = val
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
