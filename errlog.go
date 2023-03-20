package errlog

import "strings"

type exprStrings interface {
	exprStrings(dest []string)
}

func exprStringsToString(exprStrings exprStrings) string {
	target := make([]string, 0)
	exprStrings.exprStrings(target)
	return strings.Join(target, "")
}

type exprStringsFunc func(dest []string)

func (fn exprStringsFunc) exprStrings(dest []string) { fn(dest) }

type Pair struct {
	Key  string
	expr exprStrings
}

type LogType int

const (
	LogDebug LogType = 1 << iota
	LogInfo
	LogFixed
	LogWarn
	LogError
	LogFatal
	UnknownLog     LogType = 0
	LogDebugExpr           = "DEBUG"
	LogInfoExpr            = "INFO"
	LogFixedExpr           = "FIXED"
	LogWarnExpr            = "WARN"
	LogErrorExpr           = "ERROR"
	LogFatalExpr           = "FATAL"
	LogUnknownExpr         = "UNKNOWN"
)
