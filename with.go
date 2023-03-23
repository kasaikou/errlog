package errlog

import (
	"reflect"
	"strconv"

	"golang.org/x/exp/constraints"
)

func QuotedString(key, str string) Pair {
	return Pair{
		Key: key,
		expr: exprStringsFunc(func() []string {
			return []string{"'", str, "'"}
		}),
	}
}

func Int[T constraints.Signed](key string, signed T) Pair {
	return Pair{
		Key: key,
		expr: exprStringsFunc(func() []string {
			return []string{strconv.FormatInt(int64(signed), 10)}
		}),
	}
}

func Uint[T constraints.Unsigned](key string, unsigned T) Pair {
	return Pair{
		Key: key,
		expr: exprStringsFunc(func() []string {
			return []string{
				strconv.FormatUint(uint64(unsigned), 10), " (0x", strconv.FormatUint(uint64(unsigned), 16), ")",
			}
		}),
	}
}

func Bool(key string, boolean bool) Pair {
	return Pair{
		Key: key,
		expr: exprStringsFunc(func() []string {
			return []string{strconv.FormatBool(boolean)}
		}),
	}
}

func TypeOf(key string, value interface{}) Pair {
	return Pair{
		Key: key,
		expr: exprStringsFunc(func() []string {
			return []string{reflect.TypeOf(value).String()}
		}),
	}
}
