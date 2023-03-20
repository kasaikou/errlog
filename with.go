package errlog

import (
	"reflect"
	"strconv"

	"golang.org/x/exp/constraints"
)

func QuotedString(key, str string) Pair {
	return Pair{
		Key: key,
		expr: exprStringsFunc(func(dest []string) {
			dest = dest[:0]
			dest = append(dest, "'", str, "'")
			_ = dest
		}),
	}
}

func Int[T constraints.Signed](key string, signed T) Pair {
	return Pair{
		Key: key,
		expr: exprStringsFunc(func(dest []string) {
			dest = dest[:0]
			dest = append(dest, strconv.FormatInt(int64(signed), 10))
			_ = dest
		}),
	}
}

func Uint[T constraints.Unsigned](key string, unsigned T) Pair {
	return Pair{
		Key: key,
		expr: exprStringsFunc(func(dest []string) {
			dest = dest[:0]
			dest = append(dest, strconv.FormatUint(uint64(unsigned), 10), " (0x"+strconv.FormatUint(uint64(unsigned), 16)+")")
			_ = dest
		}),
	}
}

func Bool(key string, boolean bool) Pair {
	return Pair{
		Key: key,
		expr: exprStringsFunc(func(dest []string) {
			dest = dest[:0]
			dest = append(dest, strconv.FormatBool(boolean))
			_ = dest
		}),
	}
}

func TypeOf(key string, value interface{}) Pair {
	return Pair{
		Key: key,
		expr: exprStringsFunc(func(dest []string) {
			dest = dest[:0]
			dest = append(dest, reflect.TypeOf(value).String())
			_ = dest
		}),
	}
}
