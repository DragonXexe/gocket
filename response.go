package gocket

import "encoding/json"

type Response interface {
	StatusCode() int
	Content() []byte
}
type RawHTMLWithCode struct {
	Code int
	HTML string
}

type RawHTML string

func (h RawHTML) StatusCode() int {
	return 200
}

func (h RawHTML) Content() []byte {
	return []byte(h)
}

type jSONResponse[T any] struct {
	code int
	val  T
}

func (j jSONResponse[T]) StatusCode() int {
	return j.code
}

func (j jSONResponse[T]) Content() []byte {
	bytes, err := json.Marshal(j.val)
	if err != nil {
		panic(err)
	}
	return bytes
}

func JSONResponse[T any](code int, val T) Response {
	return jSONResponse[T]{
		code: code,
		val:  val,
	}
}
