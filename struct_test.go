package flags

import (
	"flag"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

// -- string Value
type stringValue string

func newStringValue(val string) *stringValue {
	return (*stringValue)(&val)
}

func (s *stringValue) Set(val string) error {
	*s = stringValue(val + " world")
	return nil
}
func (s *stringValue) String() string { return string(*s) }

func TestStruct(t *testing.T) {
	type Named struct {
		String string
	}
	newErr := fmt.Errorf
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
			name: "with_tag",
			conf: &struct {
				String string `flags:"s"`
			}{},
			args: []string{"-String=hello"},
			perr: newErr("flag provided but not defined: -String"),
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
			exp:  &struct{ Val flag.Value }{newStringValue("hello world")},
			args: []string{"-Val=hello"},
		},
		{
			name: "Named",
			conf: &Named{},
			exp:  &Named{"hello"},
			args: []string{"-Named.String=hello"},
		},
		{
			name: "nested",
			conf: &struct{ Nested Named }{},
			exp:  &struct{ Nested Named }{Named{"hello"}},
			args: []string{"-Nested.Named.String=hello"},
		},
		{
			name: "embedded",
			conf: &struct{ Named }{},
			exp:  &struct{ Named }{Named{"hello"}},
			args: []string{"-Named.String=hello"},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			fs := flag.NewFlagSet(tc.name, flag.ContinueOnError)
			fs.SetOutput(ioutil.Discard)

			err := Struct(tc.conf, fs)
			if tc.err != nil || err != nil {
				if tc.err == nil {
					assert.NoError(t, err)
				} else {
					assert.EqualError(t, err, tc.err.Error())
				}
				return
			}

			err = fs.Parse(tc.args)
			if tc.perr != nil || err != nil {
				if tc.perr == nil {
					assert.NoError(t, err)
				} else {
					assert.EqualError(t, err, tc.perr.Error())
				}
				return
			}
			assert.Equal(t, tc.exp, tc.conf)
		})
	}
}

func TestStruct_ignore(t *testing.T) {
	// {
	// 	name: "ignore",
	// 	args: []string{"-String=hello"},
	// 	perr: newErr("flag provided but not defined: -String"),
	// },

	conf := &struct {
		String string `flags:"-"`
	}{}
	fs := flag.NewFlagSet("ignore", flag.ContinueOnError)
	err := Struct(conf, fs)
	if err != nil {
		t.Errorf("got an error: %v", err)
	}
	if fs.Lookup("-") != nil {
		t.Errorf("expected not to create flag -")
	}
}

func TestStruct_usage(t *testing.T) {
	conf := &struct {
		String string `flags:",set a string"`
	}{}
	fs := flag.NewFlagSet("ignore", flag.ContinueOnError)
	err := Struct(conf, fs)
	if err != nil {
		t.Errorf("got an error: %v", err)
	}
	exp := "set a string"
	assert.Equal(t, exp, fs.Lookup("String").Usage)
}
