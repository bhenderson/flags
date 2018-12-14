package flags

import (
	"flag"
	"fmt"
	"reflect"
	"strings"
)

var ErrNotStruct = fmt.Errorf("not a pointer to a struct")

// Struct takes a pointer to a struct and sets flags based on the fields.
// Field names are prefixed by the exported (unless anonymous) name.
// The struct tag `flags:"name,usage"` is supported. "-" will ignore the field.
// Usage is optional.
func Struct(v interface{}, fs *flag.FlagSet) error {
	if v == nil {
		return nil
	}
	val := reflect.ValueOf(v)
	return addStruct(val, "", fs)
}

var typeOfFlagValue = reflect.TypeOf((*flag.Value)(nil)).Elem()

func addStruct(v reflect.Value, prefix string, fs *flag.FlagSet) error {
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return ErrNotStruct
	}
	v = v.Elem()
	vs := v.Type()
	prefix = joinPath(prefix, vs.Name())
	for i := 0; i < v.NumField(); i++ {
		sf := vs.Field(i)
		name := prefix
		if !sf.Anonymous {
			name = joinPath(name, sf.Name)
		}
		usage := ""
		tag := sf.Tag.Get("flags")
		if tag != "" {
			ts := strings.SplitN(tag, ",", 2)
			if ts[0] != "" {
				name = ts[0]
				if name == "-" {
					continue
				}
			}
			if len(ts) > 1 {
				usage = ts[1]
			}
		}
		sv := v.Field(i)
		if sv.Kind() != reflect.Ptr && !sv.Type().Implements(typeOfFlagValue) {
			sv = sv.Addr()
		}
		Flag(sv.Interface(), name, usage, fs)
	}
	return nil
}

func joinPath(a, b string) string {
	if a != "" && b != "" {
		a += "."
	}
	return a + b
}
