package errlog_test

import (
	"errors"
	"testing"

	"github.com/streamwest-1629/errlog"
)

func TestExamples(t *testing.T) {
	Example()
}

func StacktraceTestA() error {
	return errlog.Stamp(errlog.Container{
		Type:    errlog.LogError,
		Message: "call in stack",
		With:    []errlog.Pair{errlog.QuotedString("func", "StacktraceTestA")},
	})
}

func StacktraceTestB() error {

	if err := StacktraceTestA(); err != nil {
		return errlog.Stamp(errlog.Container{
			Type:     errlog.LogError,
			Message:  "call in stack",
			With:     []errlog.Pair{errlog.QuotedString("func", "StacktraceTestB")},
			Internal: err,
		})
	}

	return nil
}

func Example() {
	errlog.Default.Message(errlog.LogInfo, nil, "hello, world errlog!")
	errlog.Default.Log(errlog.Stamp(
		errlog.Container{
			Type:    errlog.LogInfo,
			Message: "errlog was created with consideration for its compatibility with errors in Go language",
			Description: "the objects issued for each logging can be treated as error types." +
				" Therefore, you can choose whether to output the created logs as \"output as 'logged message'\"" +
				" or \"delegate the processing to the caller as an 'error'\".",
		}))

	errlog.Default.Message(errlog.LogInfo, nil, "it supports errors.Unwrap(), errors.Is() function in errors package")
	errlog.Default.Message(errlog.LogInfo, errors.New("example error"), "of course, it also supports the normal errors.New() function")
	errlog.Default.Message(errlog.LogFixed, StacktraceTestB(), "it supports stacktrace")
}
