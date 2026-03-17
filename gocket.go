// Package gocket
// This is a package called gocket
package gocket

import (
	"fmt"
	"net/http"
	"sync"
)

type Gocket struct {
	// This tree is build from
	routes radixTree
	// See comment on `gocket.GocketCtx.localState`
	state    map[string]any
	contexts sync.Pool
}

func NewGocket() *Gocket {
	return &Gocket{
		routes: newRadixTree(),
		state:  make(map[string]any),
		contexts: sync.Pool{
			New: func() any {
				return GocketCtx{}
			},
		},
	}
}

func (g *Gocket) AddRoute(route Route) {
	g.routes.addRoute(route)
}

func (g *Gocket) Run(port string) {
	addr := fmt.Sprintf("0.0.0.0:%s", port)
	fmt.Printf("Mounting on %s\n", addr)
	fmt.Println("Routes:")
	for route := range g.routes.Routes() {
		fmt.Print("    ")
		fmt.Println(route.Pattern.String())
	}
	http.ListenAndServe(addr, g)
}

func (g *Gocket) ServeHTTP(responder http.ResponseWriter, rawReq *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Route panicked: %s\n", err)
			responder.WriteHeader(500)
			responder.Write([]byte("internal server error"))
			return
		}
	}()

	possiblePaths := g.routes.matchPath(splitPathToParts(rawReq.URL.Path))

	for _, route := range possiblePaths {
		req, err := route.Pattern.ParseRequest(rawReq)
		if err != nil {
			continue
		}
		ctx := g.contexts.Get().(GocketCtx)
		ctx.reset(g, &req, responder, rawReq)

		for _, middleWare := range route.middleWares {
			res := middleWare(&ctx)
			code := res.code()
			if code == middleWareSkip {
				break
			} else if code == middleWareBlock {
				response := res.reason()
				writeResponse(&ctx, response)
				return
			}
		}
		fmt.Println("Matched request")
		response := route.Handler(&ctx)
		writeResponse(&ctx, response)
		return

	}
	writeNotFound(responder)
}

func writeResponse(ctx *GocketCtx, response Response) {
	// if response is nil assume that the handler did something special
	if response == nil {
		return
	}
	ctx.writer.WriteHeader(response.StatusCode())
	ctx.writer.Write(response.Content())
}

func writeNotFound(writer http.ResponseWriter) {
	fmt.Println("Failed to match request")
	writer.WriteHeader(404)
	writer.Write([]byte("Page Not Found"))
}
