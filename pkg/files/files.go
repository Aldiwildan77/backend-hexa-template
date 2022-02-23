package files

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// FileExists to check if file by path are exists
func FileExists(path string) bool {
	i, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return !i.IsDir()
}

// Trace to get separated trace
func Trace() (string, int, string) {
	pc, file, line, ok := runtime.Caller(1)
	if ok {
		return file, line, runtime.FuncForPC(pc).Name()
	}

	return "", 0, ""
}

// PrintStackTrace to get current stack trace
func PrintStackTrace() string {
	trace := ""

	pc, file, line, ok := runtime.Caller(1)
	if ok {
		trace = fmt.Sprintf("Called from %s, line #%d, func: %v\n", file, line, runtime.FuncForPC(pc).Name())
	}

	return trace
}

// PrintFileTrace to get current file
func PrintFileTrace() string {
	_, file, _, ok := runtime.Caller(1)
	if ok {
		return file
	}

	return ""
}

// PrintLineTrace to get current line
func PrintLineTrace() int {
	_, _, line, ok := runtime.Caller(1)
	if ok {
		return line
	}

	return 0
}

// PrintFuncTrace to get current function
func PrintFuncTrace() string {
	pc, _, _, ok := runtime.Caller(1)
	if ok {
		return runtime.FuncForPC(pc).Name()
	}

	return ""
}

// PrintPackageTrace to get current package
func PrintPackageTrace() string {
	pc, _, _, ok := runtime.Caller(1)
	if ok {
		return strings.Split(runtime.FuncForPC(pc).Name(), ".")[0]
	}

	return ""
}
