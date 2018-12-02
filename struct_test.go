package flags

import (
	"flag"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
)

// -- string Value
type stringValue string

func newStringValue(val string) *stringValue {
	return (*stringValue)(&val)
}

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}
func (s *stringValue) String() string { return string(*s) }

func TestStruct(t *testing.T) {
	newErr := fmt.Errorf
	equal := func(a, b error) bool {
		if a == nil || b == nil {
			return a == b
		}
		return a.Error() == b.Error()
	}

	tcs := []struct {
		name      string
		conf, exp interface{}
		args      []string
		err       error
		perr      error
	}{
		{
			name: "set_string",
			conf: &struct{ String string }{},
			exp:  &struct{ String string }{"hello"},
			args: []string{"-String=hello"},
		},
		{
			name: "with_tag",
			conf: &struct {
				String string `flags:"s"`
			}{},
			exp: &struct {
				String string `flags:"s"`
			}{"hello"},
			args: []string{"-s=hello"},
		},
		{
			name: "string_ptr",
			conf: &struct{ String *string }{new(string)},
			exp:  &struct{ String *string }{func() *string { s := "hello"; return &s }()},
			args: []string{"-String=hello"},
		},
		{
			name: "nil_ptr",
			conf: &struct{ String *string }{},
			args: []string{"-String=hello"},
			perr: newErr("flag provided but not defined: -String"),
		},
		{
			name: "all_nils",
		},
		{
			name: "not_a_ptr",
			conf: struct{}{},
			err:  ErrNotStruct,
		},
		{
			name: "flag.Value",
			conf: &struct{ Val flag.Value }{newStringValue("")},
			exp:  &struct{ Val flag.Value }{newStringValue("hello")},
			args: []string{"-Val=hello"},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			fs := flag.NewFlagSet(tc.name, flag.ContinueOnError)
			fs.SetOutput(ioutil.Discard)

			err := Struct(tc.conf, fs)
			if tc.err != nil || err != nil {
				if !equal(tc.err, err) {
					t.Errorf("expected error(%v) got error(%v)", tc.err, err)
				}
				return
			}

			err = fs.Parse(tc.args)
			if tc.perr != nil || err != nil {
				if !equal(tc.perr, err) {
					t.Errorf("expected parse error(%v) got error(%v)", tc.perr, err)
				}
				return
			}
			if !reflect.DeepEqual(tc.exp, tc.conf) {
				t.Errorf("expected:\n%#v\ngot:\n%#v", tc.exp, tc.conf)
			}
		})
	}
}
