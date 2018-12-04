package flags

import (
	"flag"
	"fmt"
	"reflect"
	"time"
)

// Flag takes an interface and sets it on the given FlagSet
func Flag(v interface{}, name, usage string, fs *flag.FlagSet) {
	val := reflect.ValueOf(v)
	if v == nil || val.Kind() == reflect.Ptr && val.IsNil() {
		return
	}
	if fs == nil {
		fs = flag.CommandLine
	}

	switch x := v.(type) {
	case *bool:
		fs.BoolVar(x, name, *x, usage)
	case *time.Duration:
		fs.DurationVar(x, name, *x, usage)
	case *float64:
		fs.Float64Var(x, name, *x, usage)
	case *int64:
		fs.Int64Var(x, name, *x, usage)
	case *int:
		fs.IntVar(x, name, *x, usage)
	case *string:
		fs.StringVar(x, name, *x, usage)
	case *uint64:
		fs.Uint64Var(x, name, *x, usage)
	case *uint:
		fs.UintVar(x, name, *x, usage)
	case flag.Value:
		fs.Var(x, name, usage)
	default:
		panic(fmt.Sprintf("unsupported type %T", v))
	}
}
