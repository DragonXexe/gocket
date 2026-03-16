package gocket

import "net/http"

func (ctx *GocketCtx) GetCookie(name string) (string, bool) {
	for _, cookie := range ctx.Req.Cookies {
		if cookie.Name == name {
			return cookie.Value, true
		}
	}
	return "", false
}
func (ctx *GocketCtx) SetCookie(cookie http.Cookie) {
	ctx.writer.Header().Add("Set-Cookie", cookie.String())
}
