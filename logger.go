package gocket

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func DebugPrint(v any) {
	_, file, line, _ := runtime.Caller(1)
	file = filepath.Base(file)
	fmt.Printf("[%s:%d] %#v\n", file, line, v)
}

func LogPrint(s string, a ...any) {
	message := fmt.Sprintf(s, a...)
	_, file, _, _ := runtime.Caller(1)
	pkg := filepath.Base(filepath.Dir(file))
	base := filepath.Base(file)

	fmt.Printf("[%s/%s] %s\n", pkg, base, message)
}
