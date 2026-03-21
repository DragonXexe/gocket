package gocket

import (
	"fmt"
	"path/filepath"
	"runtime"
)

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	green  = "\033[32m"
)

func DebugPrint(v any) {
	_, file, line, _ := runtime.Caller(1)
	file = filepath.Base(file)
	fmt.Printf("[%s:%d] %#v\n", file, line, v)
}

func LogInfo(s string, a ...any) {
	message := fmt.Sprintf(s, a...)
	_, file, _, _ := runtime.Caller(1)
	pkg := filepath.Base(filepath.Dir(file))
	base := filepath.Base(file)

	fmt.Printf("%s[%s/%s]%s %s\n", blue, pkg, base, reset, message)
}
func LogWarning(s string, a ...any) {
	message := fmt.Sprintf(s, a...)
	_, file, _, _ := runtime.Caller(1)
	pkg := filepath.Base(filepath.Dir(file))
	base := filepath.Base(file)

	fmt.Printf("%s[%s/%s]%s %s\n", yellow, pkg, base, reset, message)
}
func LogError(s string, a ...any) {
	message := fmt.Sprintf(s, a...)
	_, file, _, _ := runtime.Caller(1)
	pkg := filepath.Base(filepath.Dir(file))
	base := filepath.Base(file)

	fmt.Printf("%s[%s/%s]%s %s\n", red, pkg, base, reset, message)
}
