package flags

import (
	"flag"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestEnv(t *testing.T) {
	tcs := []struct {
		env      []string
		expected string
	}{
		{
			[]string{
				"A-STRIN=ignore",
				"A-STRING=ohmy",
				"A-STRINGG=ignore",
			},
			"ohmy",
		},
	}

	for i, tc := range tcs {
		name := strconv.Itoa(i)
		t.Run(name, func(t *testing.T) {
			defer setEnv(tc.env)()

			fs := flag.NewFlagSet(name, flag.ContinueOnError)
			as := fs.String("a-string", "default", "")

			Env(fs)
			if *as != tc.expected {
				t.Errorf("expected -a-string=%s got %q", tc.expected, *as)
			}
		})
	}
}

func setEnv(env []string) func() {
	for _, e := range env {
		ss := strings.SplitN(e, "=", 2)
		os.Setenv(ss[0], ss[1])
	}
	return func() {
		for _, e := range env {
			ss := strings.SplitN(e, "=", 2)
			os.Unsetenv(ss[0])
		}
	}
}
