package errlog

import (
	"fmt"
	"io"
	"os"
	"strings"
)

var LoggingTypes = LogWarn | LogError | LogFatal

var Default interface {
	Message(level LogType, err error, msg string, with ...Pair)
	Log(content Container)
} = DefaultLoggers.CommandLine

var DefaultLoggers = struct {
	CommandLine CommandLineLogger
}{
	CommandLine: CommandLineLogger{
		Color:        true,
		Dest:         os.Stdout,
		DefaultStamp: DefaultStamp,
	},
}

func writeStrings(dest io.Writer, strs ...string) (int, error) {
	n := 0
	for _, str := range strs {
		if w, e := dest.Write([]byte(str)); e != nil {
			return n, e
		} else {
			n += w
		}
	}

	return n, nil
}

type CommandLineLogger struct {
	Color        bool
	Dest         io.Writer
	DefaultStamp LogStampConfig
}

func (cl CommandLineLogger) Message(level LogType, err error, msg string, with ...Pair) {
	cl.Log(cl.DefaultStamp.StampDetail(Container{
		Type:     level,
		Message:  msg,
		Internal: err,
	}, 1))
}

func (cl CommandLineLogger) Log(content Container) {

	content = cl.DefaultStamp.StampDetail(content, 1)
	const (
		FgReset   = "\x1b[0m"
		FgRed     = "\x1b[31m"
		FgGreen   = "\x1b[32m"
		FgYellow  = "\x1b[33m"
		FgBlue    = "\x1b[34m"
		FgMagenta = "\x1b[35m"
		FgCyan    = "\x1b[36m"
		FgWhite   = "\x1b[37m"
	)

	checkAndPrint := func(color string, strs ...string) (err error) {

		err = nil

		if cl.Color {
			if _, writeErr := writeStrings(cl.Dest, color); writeErr != nil {
				err = writeErr
			} else if _, writeErr := writeStrings(cl.Dest, strs...); writeErr != nil {
				err = writeErr
			} else if _, writeErr := writeStrings(cl.Dest, FgReset); writeErr != nil {
				err = writeErr
			}
		} else {
			if _, writeErr := writeStrings(cl.Dest, strs...); writeErr != nil {
				err = writeErr
			}
		}

		if err != nil {
			return fmt.Errorf("cannot write strings: %w", err)
		}
		return nil
	}

	withCache := make([]string, 256)
	content.WalkErrorstack(func(container Container, wrapped error) {

		if container.stampConfig.EnableTimestamp {
			checkAndPrint(FgCyan, container.timestamp.Format(container.stampConfig.TimestampFormat), " ")
		}

		type TypeColorExpr struct {
			Type  LogType
			Color string
			Expr  string
		}

		typeColors := [...]TypeColorExpr{
			{Type: LogDebug, Color: FgBlue, Expr: LogDebugExpr},
			{Type: LogInfo, Color: FgGreen, Expr: LogInfoExpr},
			{Type: LogFixed, Color: FgCyan, Expr: LogFixedExpr},
			{Type: LogWarn, Color: FgYellow, Expr: LogWarnExpr},
			{Type: LogError, Color: FgMagenta, Expr: LogErrorExpr},
			{Type: LogFatal, Color: FgRed, Expr: LogFatalExpr},
		}

		msgColor := FgReset
		checkAndPrint(FgReset, "[")
		for _, typeColor := range typeColors {
			if container.Type&typeColor.Type != 0 {
				if msgColor != FgReset {
					checkAndPrint(FgReset, "|")
				}
				msgColor = typeColor.Color
				checkAndPrint(msgColor, typeColor.Expr)
			}
		}
		checkAndPrint(FgReset, "] ")
		checkAndPrint(msgColor, container.Message)

		if len(container.With) > 0 {
			checkAndPrint(FgReset, " (")
			checkAndPrint(FgCyan, container.With[0].Key)
			checkAndPrint(FgReset, ": ")
			container.With[0].expr.exprStrings(withCache)
			checkAndPrint(msgColor, withCache...)
			for _, with := range container.With[1:] {
				checkAndPrint(FgReset, ", ")
				checkAndPrint(FgCyan, with.Key)
				checkAndPrint(FgReset, ": ")
				container.With[0].expr.exprStrings(withCache)
				checkAndPrint(msgColor, withCache...)
			}
			checkAndPrint(FgReset, ")")
		}

		if wrapped != nil {
			checkAndPrint(FgReset, ": ")
			checkAndPrint(FgMagenta, wrapped.Error())
		}

		checkAndPrint(FgReset, "\n")
		if container.Description != "" {
			checkAndPrint(FgReset, "   description: ", strings.TrimSuffix(container.Description, "\n"), "\n")
		}

	})

	if content.stackTraces != nil {
		n := content.StackTrace(withCache)
		if n > 0 {
			withCache := withCache[:n]
			checkAndPrint(FgReset, "stacktrace:\n")
			for _, stack := range withCache {
				checkAndPrint(FgReset, "  ", stack, "\n")
			}
		}
	}
}
