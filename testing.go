package gocket

import (
	"context"
	"fmt"
	"io"
	"net/http/httptest"
)

type TestConfig struct {
	Gocket      *Gocket
	Method      string
	Pattern     string
	Path        string
	Body        io.Reader
	Ctx         context.Context
	Cancel      context.CancelFunc
	LocalState  map[string]any
	MiddleWares []MiddleWare
	Handler     Handler
}

func TestHandler(config *TestConfig, handler Handler) (resp Response, err error) {
	ctx, r := TestCreateContext(config)
	if r != nil {
		err = r
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			resp = nil
			err = fmt.Errorf("Handler panicked: %v", r)
		}
	}()
	// run middlewares
	for _, middleware := range config.MiddleWares {
		res := middleware(&ctx)
		if res.code() == middleWareBlock {
			return res.reason(), nil
		} else if res.code() == middleWarePass {
			continue
		} else if res.code() == middleWareSkip {
			return nil, fmt.Errorf("failed to test handler because a middleware resulted in skip action")
		}
	}
	resp = handler(&ctx)
	return resp, nil
}

func TestCreateContext(config *TestConfig) (GocketCtx, error) {
	rec := httptest.NewRecorder()
	raw := httptest.NewRequest(config.Method, config.Path, config.Body)
	if config.Ctx == nil {
		ctx, cancel := context.WithCancel(raw.Context())
		config.Ctx = ctx
		config.Cancel = cancel
	}
	pattern, err := ParsePattern(config.Method, config.Pattern)
	if err != nil {
		return GocketCtx{}, err
	}
	req, err := pattern.ParseRequest(raw)
	if err != nil {
		return GocketCtx{}, err
	}
	if config.Gocket == nil {
		config.Gocket = &Gocket{}
	}
	return GocketCtx{
		gocket:        config.Gocket,
		context:       config.Ctx,
		cancel:        config.Cancel,
		Req:           &req,
		writer:        rec,
		origalRequest: raw,
		localState:    config.LocalState,
	}, nil
}
