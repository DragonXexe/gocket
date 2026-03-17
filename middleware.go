package gocket

type (
	MiddleWare func(ctx *GocketCtx) MiddleWareResult
)

func (group *Group) AddMiddleWare(middleWare MiddleWare) {
	group.middleWares = append(group.middleWares, middleWare)
}

type middleWareResultCode uint

const (
	middleWarePass middleWareResultCode = iota
	middleWareSkip
	middleWareBlock
)

type MiddleWareResult interface {
	code() middleWareResultCode
	// may return nil if their is no reason
	reason() Response
}

type blocked struct {
	response Response
}

func (b blocked) code() middleWareResultCode {
	return middleWareBlock
}

func (b blocked) reason() Response {
	return b.response
}

type pass struct{}

func (p pass) code() middleWareResultCode {
	return middleWarePass
}

func (p pass) reason() Response {
	return nil
}

type skip struct{}

func (s skip) code() middleWareResultCode {
	return middleWareSkip
}

func (s skip) reason() Response {
	return nil
}

func Pass() MiddleWareResult {
	return pass{}
}

func Skip() MiddleWareResult {
	return skip{}
}

func Block(response Response) MiddleWareResult {
	return blocked{response}
}
