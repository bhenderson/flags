package flags

import (
	"flag"
	"fmt"
	"reflect"
	"time"
)

// Flag takes an interface and sets it on the given FlagSet
// It panics if v is an unsupported type. It takes everything supported by the
// flag package, as well as structs using Struct
func Flag(v interface{}, name, usage string, fs *flag.FlagSet) {
	// log.Printf("adding flag %#v (%s)", v, name)
	if v == nil {
		return
	}
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr && val.IsNil() {
		return
	}
	if fs == nil {
		fs = flag.CommandLine
	}

	if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
		addStruct(val, name, fs)
		return
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
