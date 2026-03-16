package gocket

import (
	"github.com/coder/websocket"
)

type WebSockHandler func(con *websocket.Conn)

func (ctx *GocketCtx) UpgradeWebsocket(handler WebSockHandler) error {
	con, err := websocket.Accept(ctx.writer, ctx.origalRequest, nil)
	if err != nil {
		return err
	}
	go func() {
		// make sure that the connection is closed no matter what handler does
		defer con.CloseNow()
		handler(con)
	}()

	return nil
}
