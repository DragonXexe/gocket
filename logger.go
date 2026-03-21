package gocket

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
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

func printLog(color string, message string) {
	_, file, _, _ := runtime.Caller(2)

	pkg := filepath.Base(filepath.Dir(file))
	base := filepath.Base(file)

	// Strip the @v0.0.0-... version suffix from the package dir.
	if i := strings.Index(pkg, "@"); i != -1 {
		pkg = pkg[:i]
	}
	fmt.Printf("%s[%s/%s]%s %s\n", color, pkg, base, reset, message)
}

func LogSuccess(s string, a ...any) {
	message := fmt.Sprintf(s, a...)
	printLog(green, message)
}

func LogInfo(s string, a ...any) {
	message := fmt.Sprintf(s, a...)
	printLog(blue, message)
}

func LogWarning(s string, a ...any) {
	message := fmt.Sprintf(s, a...)
	printLog(yellow, message)
}

func LogError(err error) {
	printLog(red, err.Error())
}

func LogErrorf(s string, a ...any) {
	message := fmt.Sprintf(s, a...)
	printLog(red, message)
}
