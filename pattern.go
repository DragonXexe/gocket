package gocket

import (
	"fmt"
	"slices"
	"strings"
)

// Pattern ...
// Examples:
// "/path/to/{dynamic_path}?{query1}&{query2}"
type Pattern struct {
	Method  string
	Path    Path
	Queries []string
}

func (pattern Pattern) String() string {
	var s strings.Builder
	s.WriteString(pattern.Method)
	s.WriteString(" ")
	for _, part := range pattern.Path.Parts {
		if part.Wildcard {
			fmt.Fprintf(&s, "/{%s}", part.Name)
		} else {
			fmt.Fprintf(&s, "/%s", part.Name)
		}
	}
	if len(pattern.Queries) != 0 {
		s.WriteString("?")
	}
	for i, query := range pattern.Queries {
		fmt.Fprintf(&s, "{%s}", query)
		if i+1 < len(pattern.Queries) {
			s.WriteString("&")
		}
	}
	return s.String()
}

func ParsePattern(method string, pattern string) (Pattern, error) {
	path, queries, hasQueries := strings.Cut(pattern, "?")
	parsedPath, err := ParsePath(path)
	if err != nil {
		return Pattern{}, err
	}
	if !hasQueries {
		return Pattern{Method: method, Path: parsedPath}, nil
	}
	foundQueries := make([]string, 0)
	for queryOriginal := range strings.SplitSeq(queries, "&") {
		query, foundPrefix := strings.CutPrefix(queryOriginal, "{")
		if !foundPrefix {
			return Pattern{}, fmt.Errorf("invalid query paramter: %s", queryOriginal)
		}
		query, foundPostfix := strings.CutSuffix(query, "}")
		if !foundPostfix {
			return Pattern{}, fmt.Errorf("invalid query paramter: %s", queryOriginal)
		}
		if slices.Contains(foundQueries, query) {
			return Pattern{}, fmt.Errorf("duplicate query paramter: %s", query)
		}
		foundQueries = append(foundQueries, query)
	}
	return Pattern{
		Method:  method,
		Path:    parsedPath,
		Queries: foundQueries,
	}, nil
}

type Path struct {
	Parts     []PathElement
	HasAnyEnd bool
}
type PathElement struct {
	Wildcard bool
	Name     string
}

func ParsePath(path string) (Path, error) {
	// remove leading and trailing slash
	path, _ = strings.CutPrefix(path, "/")
	path, _ = strings.CutSuffix(path, "/")
	if path == "" {
		return Path{[]PathElement{}, false}, nil
	}
	parts := make([]PathElement, 0)
	for part := range strings.SplitSeq(path, "/") {
		if part == "" {
			return Path{}, fmt.Errorf("invalid path element, got an empty string")
		}
		if !strings.Contains(part, "{") && !strings.Contains(part, "}") {
			parts = append(parts, PathElement{Wildcard: false, Name: part})
			continue
		}
		originalPart := part
		part, foundPrefix := strings.CutPrefix(originalPart, "{")
		if !foundPrefix {
			return Path{}, fmt.Errorf("invalid path element: %s", originalPart)
		}
		part, foundPostfix := strings.CutSuffix(part, "}")
		if !foundPostfix {
			return Path{}, fmt.Errorf("invalid path element: %s", originalPart)
		}
		element := PathElement{Wildcard: true, Name: part}
		if slices.Contains(parts, element) {
			return Path{}, fmt.Errorf("duplicate path wildcard: %s", part)
		}
		parts = append(parts, element)
	}
	return Path{parts, false}, nil
}
