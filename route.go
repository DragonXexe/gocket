package gocket

type Handler func(ctx *GocketCtx) Response

type Route struct {
	middleWares []MiddleWare
	Pattern     Pattern
	Handler     Handler
}
