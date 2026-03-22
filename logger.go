package gocket

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"
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

func printLog(color string, level string, message string) {
	_, file, _, _ := runtime.Caller(2)

	pkg := filepath.Base(filepath.Dir(file))

	// Strip the @v0.0.0-... version suffix from the package dir.
	if i := strings.Index(pkg, "@"); i != -1 {
		pkg = pkg[:i]
	}
	timeString := time.Now().Format("15:04:05")
	fmt.Printf("[%s] %s[%s]%s [%s] %s\n", timeString, color, level, reset, pkg, message)
}

func LogSuccess(s string, a ...any) {
	message := fmt.Sprintf(s, a...)
	printLog(green, "SUC", message)
}

func LogInfo(s string, a ...any) {
	message := fmt.Sprintf(s, a...)
	printLog(blue, "INF", message)
}

func LogWarning(s string, a ...any) {
	message := fmt.Sprintf(s, a...)
	printLog(yellow, "WRN", message)
}

func LogError(err error) {
	printLog(red, "ERR", err.Error())
}

func LogErrorf(s string, a ...any) {
	message := fmt.Sprintf(s, a...)
	printLog(red, "ERR", message)
}
