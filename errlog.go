package errlog

import "strings"

type exprStrings interface {
	exprStrings() []string
}

func exprStringsToString(exprStrings exprStrings) string {
	return strings.Join(exprStrings.exprStrings(), "")
}

type exprStringsFunc func() []string

func (fn exprStringsFunc) exprStrings() []string { return fn() }

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
