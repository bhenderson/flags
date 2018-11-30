package flags

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	defer writeConfFile(t)()

	tcs := []struct {
		args     []string
		expected string
	}{
		{
			[]string{
				"-a-string", "a",
			},
			"a",
		},
		{
			[]string{
				"-config", confFile,
			},
			"hello",
		},
		{
			[]string{
				"-config", confFile, "-a-string=b",
			},
			"b",
		},
		{
			[]string{
				"-config", confFile, "-a-string=b", "-config", confFile,
			},
			"hello",
		},
	}

	for _, tc := range tcs {
		fs := flag.NewFlagSet(t.Name(), flag.ContinueOnError)
		as := fs.String("a-string", "default", "")
		Config("config", "", fs)

		fs.Parse(tc.args)
		if *as != tc.expected {
			t.Errorf("%v\n\texpected -a-string=%s got %q", tc.args, tc.expected, *as)
		}
	}
}

var confFile = "config.conf"

func writeConfFile(t *testing.T) func() {
	content :=
		`# a comment followed by a space

-a-string=hello
-config=` + confFile + `
extras`
	err := ioutil.WriteFile(confFile, []byte(content), 0600)
	if err != nil {
		t.Error("could not create config file", err)
	}
	return func() {
		if err := os.Remove(confFile); err != nil {
			t.Log("failed to remove test config file", err)
		}
	}
}
