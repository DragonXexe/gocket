package gocket

import (
	"strings"
)

func (g *Gocket) NewGroup(prefix string) *Group {
	prefix, _ = strings.CutPrefix(prefix, "/")
	prefix, _ = strings.CutSuffix(prefix, "/")
	return &Group{
		gocket: g,
		prefix: strings.Split(prefix, "/"),
	}
}

type Group struct {
	gocket      *Gocket
	prefix      []string
	middleWares []MiddleWare
}

func (group *Group) Handle(method string, pattern string, handler Handler) {
	parsedPattern, err := ParsePattern(method, pattern)
	if err != nil {
		panic(err)
	}
	for _, prefix := range group.prefix {
		parsedPattern.Path.Parts = insert(parsedPattern.Path.Parts, 0, PathElement{
			Wildcard: false,
			Name:     prefix,
		})
	}
	route := Route{
		Pattern:     parsedPattern,
		Handler:     handler,
		middleWares: group.middleWares,
	}
	group.gocket.AddRoute(route)
}

func (group *Group) HandleWithMiddleWare(method string, pattern string, middleWares []MiddleWare, handler Handler) {
	parsedPattern, err := ParsePattern(method, pattern)
	if err != nil {
		panic(err)
	}
	for _, prefix := range group.prefix {
		parsedPattern.Path.Parts = insert(parsedPattern.Path.Parts, 0, PathElement{
			Wildcard: false,
			Name:     prefix,
		})
	}
	totalMiddleWares := make([]MiddleWare, 0, len(group.middleWares)+len(middleWares))
	totalMiddleWares = append(totalMiddleWares, group.middleWares...)
	totalMiddleWares = append(totalMiddleWares, middleWares...)
	route := Route{
		Pattern:     parsedPattern,
		Handler:     handler,
		middleWares: totalMiddleWares,
	}
	group.gocket.AddRoute(route)
}

func insert[T any](slice []T, index int, element T) []T {
	slice = append(slice, element)
	copy(slice[index+1:], slice[index:])
	slice[index] = element
	return slice
}
