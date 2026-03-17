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
	g.state[name] = &SafeState[any]{sync.Mutex{}, value}
}

// State returns nil if the state is not present
func (ctx *GocketCtx) State(name string) *SafeState[any] {
	state, ok := ctx.localState[name]
	if ok {
		return state.(*SafeState[any])
	}
	return ctx.gocket.state[name].(*SafeState[any])
}

func (ctx *GocketCtx) anyState(name string) any {
	state, ok := ctx.localState[name]
	if ok {
		return state
	}
	return ctx.gocket.state[name]
}

func (ctx *GocketCtx) SetLocalState(name string, val any) {
	ctx.localState[name] = &SafeState[any]{sync.Mutex{}, val}
}

func GetState[T any](ctx *GocketCtx, state string) *SafeState[T] {
	s := ctx.anyState(state)
	if s == nil {
		panic(fmt.Sprintf("State was not present: %s", state))
	}
	val, ok := s.(*SafeState[T])
	if !ok {
		// This variable exists for printing T
		var t T
		panic(fmt.Sprintf("State \"%s\" was not of type: %T but of type: %T", state, t, s))
	}
	return val
}
