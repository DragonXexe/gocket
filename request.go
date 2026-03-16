package gocket

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Request struct {
	Header      http.Header
	Body        []byte
	PathValues  map[string]string
	QueryParams map[string]string
	Cookies     []*http.Cookie
}

func converHTTPRequest(req *http.Request) (Request, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return Request{}, err
	}
	return Request{
		Header: req.Header,
		Body:   body,
	}, nil
}

func JSONBody[T any](req *Request) (T, error) {
	var res T
	err := json.Unmarshal(req.Body, &res)
	if err != nil {
		return *new(T), err
	}
	return res, nil
}

func (pattern *Pattern) ParseRequest(req *http.Request) (Request, error) {
	if req.Method != pattern.Method {
		err := fmt.Errorf("mismatched method expected %s but got %s", pattern.Method, req.Method)
		return Request{}, err
	}
	pathValues, err := pattern.parsePath(req.URL.Path)
	if err != nil {
		return Request{}, err
	}
	queries, err := pattern.parseQueries(req.URL.Query())
	if err != nil {
		return Request{}, err
	}
	parsedRequest, err := converHTTPRequest(req)
	if err != nil {
		return Request{}, err
	}
	parsedRequest.PathValues = pathValues
	parsedRequest.QueryParams = queries
	parsedRequest.Cookies = req.Cookies()
	return parsedRequest, nil
}

func splitPathToParts(path string) []string {
	path, _ = strings.CutPrefix(path, "/")
	path, _ = strings.CutSuffix(path, "/")
	parts := strings.Split(path, "/")
	return parts
}

func (pattern *Pattern) parsePath(got string) (map[string]string, error) {
	parts := splitPathToParts(got)
	if len(parts) != len(pattern.Path.Parts) {
		err := fmt.Errorf("failed to parse request because length of path did not match. Expected %d but fot %d", len(pattern.Path.Parts), len(parts))
		return map[string]string{}, err
	}
	pathValues := make(map[string]string)
	for i, part := range parts {
		expected := pattern.Path.Parts[i]
		if expected.Wildcard {
			pathValues[expected.Name] = part
			continue
		}
		if expected.Name != part {
			err := fmt.Errorf("failed to parse request because a part did not match. Expected %s but fot %s", expected.Name, part)
			return map[string]string{}, err
		}
	}
	return pathValues, nil
}

func (pattern *Pattern) parseQueries(got url.Values) (map[string]string, error) {
	foundQueries := make(map[string]string)
	for _, query := range pattern.Queries {
		gotten := got.Get(query)
		if gotten == "" {
			err := fmt.Errorf("failed to parse query because %s was not set", query)
			return map[string]string{}, err
		} else {
			foundQueries[query] = gotten
		}
	}
	return foundQueries, nil
}
