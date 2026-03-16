package gocket

type (
	MiddleWare func(ctx *GocketCtx) MiddleWareResult
)

func (group *Group) AddMiddleWare(middleWare MiddleWare) {
	group.middleWares = append(group.middleWares, middleWare)
}

type MiddleWareResult interface {
	IsPass() bool
	Blocked() (Response, bool)
}
type blocked struct {
	response Response
}

func (b blocked) Blocked() (Response, bool) {
	return b.response, true
}

func (b blocked) IsPass() bool {
	return false
}

type pass struct{}

func (p pass) Blocked() (Response, bool) {
	return nil, false
}

func (p pass) IsPass() bool {
	return true
}

func Pass() MiddleWareResult {
	return pass{}
}

func Block(response Response) MiddleWareResult {
	return blocked{response}
}
