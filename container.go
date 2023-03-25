package errlog

import (
	"errors"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type LogStampConfig struct {
	EnableTimestamp       bool
	TimestampFormat       string
	EnableStackTraceTypes LogType
	StackTraceFactory     func(programCounters []uintptr, dest []string) int
}

var DefaultStamp = LogStampConfig{
	EnableTimestamp:       true,
	TimestampFormat:       time.RFC3339,
	EnableStackTraceTypes: LogFixed | LogError | LogFatal,
	StackTraceFactory: func(programCounters []uintptr, dest []string) int {

		dest = dest[:0]
		frames := runtime.CallersFrames(programCounters)

		for {
			frame, more := frames.Next()
			funcName := filepath.Base(frame.Func.Name())
			funcName = funcName[strings.IndexByte(funcName, '.')+1:]
			dest = append(dest, frame.File+":"+strconv.Itoa(frame.Line)+" "+funcName)

			if !more {
				return len(dest)
			}
		}
	},
}

func (conf *LogStampConfig) Stamp(container Container) Container {
	return conf.StampDetail(container, 1)
}

func (conf *LogStampConfig) StampDetail(container Container, traceSkip int) Container {

	if conf == nil {
		return DefaultStamp.StampDetail(container, traceSkip+1)
	}
	if container.stampConfig != nil {
		return container
	}

	if conf.EnableTimestamp {
		container.timestamp = time.Now()
	}

	if container.Type&conf.EnableStackTraceTypes != 0 {
		container.stackTraces = make([]uintptr, 64)
		nCallers := runtime.Callers(2+traceSkip, container.stackTraces)
		if nCallers == 0 {
			panic("cannot stack trace: nCallers == 0")
		}
	}

	container.stampConfig = conf
	return container
}

func Stamp(container Container) Container {
	return DefaultStamp.StampDetail(container, 1)
}

func StampDetail(container Container, traceSkip int) Container {
	return DefaultStamp.StampDetail(container, traceSkip+1)
}

type Container struct {
	stackTraces []uintptr
	timestamp   time.Time
	stampConfig *LogStampConfig
	Type        LogType
	Message     string
	Description string
	Internal    error
	With        []Pair
}

func (lc Container) Error() string {
	desc := lc.Message
	if len(lc.With) > 0 {
		desc += "("
		desc += lc.With[0].Key + ": " + exprStringsToString(lc.With[0].expr)
		for _, with := range lc.With[1:] {
			desc += ", " + with.Key + ": " + exprStringsToString(with.expr)
		}
		desc += ")"
	}

	if lc.Internal != nil {
		desc += lc.Internal.Error()
	}

	return desc
}

func (lc Container) unwrapContainer() (Container, bool) {
	if lc.Internal == nil {
		return Container{}, false
	} else if container, ok := lc.Internal.(Container); ok {
		return container, true
	}

	unwrapped := lc.Internal
	for {
		if unwrapped = errors.Unwrap(unwrapped); unwrapped == nil {
			return Container{}, false
		} else if container, ok := unwrapped.(Container); ok {
			return container, true
		}
	}
}

func (lc Container) Unwrap() error {
	if lc.Internal == nil {
		return nil
	} else if container, ok := lc.Internal.(Container); ok {
		return container.Unwrap()
	} else {
		return lc.Internal
	}
}

func (lc Container) Is(compare error) bool {
	if compare, ok := compare.(Container); ok {
		if lc.Message == compare.Message {
			return true
		} else if wrapped, ok := lc.unwrapContainer(); ok {
			return wrapped.Is(compare)
		} else {
			return false
		}

	}

	if err := lc.Unwrap(); err != nil {
		return errors.Is(err, compare)
	} else if wrapped, ok := lc.unwrapContainer(); ok {
		return wrapped.Is(err)
	} else {
		return false
	}
}

func (lc Container) StackTrace(dest []string) int {
	if internal, exist := lc.unwrapContainer(); exist {
		return internal.StackTrace(dest)
	} else if lc.stackTraces != nil && lc.stampConfig.StackTraceFactory != nil {
		return lc.stampConfig.StackTraceFactory(lc.stackTraces, dest)
	} else {
		return 0
	}
}

func (lc Container) WalkErrorstack(fn func(container Container, wrapped error)) {
	if internal, exists := lc.unwrapContainer(); exists {
		internal.WalkErrorstack(fn)
	}

	if lc.Internal != nil {
		if _, isStacked := lc.Internal.(Container); !isStacked {
			fn(lc, lc.Internal)
			return
		}
	}

	fn(lc, nil)
}
