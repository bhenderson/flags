package flags

import (
	"flag"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestStruct(t *testing.T) {
	type Config struct {
		AString string
		BString string `flags:"b,set b"`
		CString *string
	}

	makeErr := fmt.Errorf
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
		parseErr  error
	}{
		{
			"set string",
			&Config{},
			&Config{AString: "hello"},
			[]string{"-Config.AString=hello"},
			nil, nil,
		},
		{
			"with flag",
			&Config{},
			&Config{BString: "hello"},
			[]string{"-b=hello"},
			nil, nil,
		},
		{
			"string ptr",
			&Config{CString: new(string)},
			&Config{CString: func() *string { s := "hello"; return &s }()},
			[]string{"-Config.CString=hello"},
			nil, nil,
		},
		{
			"nil ptr",
			&Config{},
			&Config{},
			[]string{"-Config.CString=hello"},
			nil, makeErr("flag provided but not defined: -Config.CString"),
		},
		{
			name: "all nils",
		},
		{
			"not a ptr",
			Config{}, nil, nil, ErrNotStruct, nil,
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
			if tc.parseErr != nil || err != nil {
				if !equal(tc.parseErr, err) {
					t.Errorf("expected parse error(%v) got error(%v)", tc.parseErr, err)
				}
				return
			}
			if !reflect.DeepEqual(tc.exp, tc.conf) {
				t.Errorf("expected:\n%#v\ngot:\n%v", tc.exp, tc.conf)
			}
		})
	}
}
