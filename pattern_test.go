package gocket_test

import (
	"fmt"
	"testing"

	"github.com/DragonXexe/gocket"
)

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func TestParsePattern(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		pattern string
		want    gocket.Pattern
		wantErr bool
	}{
		{
			"simpleTest",
			"/path/to/file.html",
			gocket.Pattern{Path: must(gocket.ParsePath("/path/to/file.html")), Queries: make([]string, 0)},
			false,
		},
		{
			"queryTest",
			"/path/to/file.html?{query}",
			gocket.Pattern{Path: must(gocket.ParsePath("/path/to/file.html")), Queries: []string{"query"}},
			false,
		},
		{
			"mutlipleQueriesTest",
			"/path/to/file.html?{query1}&{query2}&{query3}",
			gocket.Pattern{Path: must(gocket.ParsePath("/path/to/file.html")), Queries: []string{"query1", "query2", "query3"}},
			false,
		},
		{
			"invalidQueryTest",
			"/path/to/file.html?{query1&query3}",
			gocket.Pattern{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := gocket.ParsePattern("GET", tt.pattern)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ParsePattern() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ParsePattern() succeeded unexpectedly")
			}
			// ignore the method for these tests
			tt.want.Method = ""
			got.Method = ""
			if fmt.Sprintf("%v", got) != fmt.Sprintf("%v", tt.want) {
				t.Errorf("ParsePattern() = %v, want %v", got, tt.want)
			}
		})
	}
}
