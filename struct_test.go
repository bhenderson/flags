package flags

import (
	"flag"
	"reflect"
	"strconv"
	"testing"
)

func TestStruct(t *testing.T) {
	type Config struct {
		AString string
		BString string `flags:"b,set b"`
		CString *string
	}

	tcs := []struct {
		conf, exp interface{}
		args      []string
		err       error
	}{
		{
			&Config{},
			&Config{AString: "hello"},
			[]string{"-Config.AString=hello"},
			nil,
		},
		{
			&Config{},
			&Config{BString: "hello"},
			[]string{"-b=hello"},
			nil,
		},
		{
			&Config{CString: new(string)},
			&Config{CString: func() *string { s := "hello"; return &s }()},
			[]string{"-Config.CString=hello"},
			nil,
		},
		{
			nil, nil, nil, nil,
		},
		{
			Config{}, nil, nil, ErrNotStruct,
		},
	}

	for i, tc := range tcs {
		name := strconv.Itoa(i)
		t.Run(name, func(t *testing.T) {
			fs := flag.NewFlagSet(name, flag.ContinueOnError)
			err := Struct(tc.conf, fs)
			if tc.err != nil {
				if tc.err != err {
					t.Errorf("expected error(%v) got error(%v)", tc.err, err)
				}
				return
			}
			fs.Parse(tc.args)
			if !reflect.DeepEqual(tc.exp, tc.conf) {
				t.Errorf("expected:\n%#v\ngot:\n%v", tc.exp, tc.conf)
			}
		})
	}
}
