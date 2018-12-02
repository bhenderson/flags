package flags

import (
	"flag"
	"fmt"
	"reflect"
	"strings"
)

var ErrNotStruct = fmt.Errorf("not a pointer to a struct")

func Struct(v interface{}, fs *flag.FlagSet) error {
	if v == nil {
		return nil
	}
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return ErrNotStruct
	}
	addStruct(val.Elem(), fs)
	return nil
}

var typeOfFlagValue = reflect.TypeOf((*flag.Value)(nil)).Elem()

func addStruct(v reflect.Value, fs *flag.FlagSet) {
	vs := v.Type()
	for i := 0; i < v.NumField(); i++ {
		sf := vs.Field(i)
		name := sf.Name
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
		if sv.IsNil() {
			return
		}
		Flag(sv.Interface(), name, usage, fs)
	}
}
