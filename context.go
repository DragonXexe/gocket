package gocket

import (
	"context"
	"net/http"
	"time"
)

type GocketCtx struct {
	gocket        *Gocket
	context       context.Context
	cancel        context.CancelFunc
	Req           *Request
	writer        http.ResponseWriter
	origalRequest *http.Request
}
func (ctx *GocketCtx) Context() context.Context {
	return ctx.context
}
func (ctx *GocketCtx) reset(g *Gocket, req *Request, writer http.ResponseWriter, originalRequest *http.Request) {
	context, cancel := context.WithTimeout(originalRequest.Context(), time.Minute)
	ctx.gocket = g
	ctx.context = context
	ctx.cancel = cancel
	ctx.Req = req
	ctx.writer = writer
	ctx.origalRequest = originalRequest.WithContext(context)
}
