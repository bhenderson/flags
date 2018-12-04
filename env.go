package flags

import (
	"flag"
	"os"
	"strings"
)

// Env parses the Environment and assigns the matching (UPPERCASED) flags
func Env(fs *flag.FlagSet) {
	if fs == nil {
		fs = flag.CommandLine
	}
	fs.VisitAll(func(f *flag.Flag) {
		name := strings.ToUpper(f.Name)
		if v, ok := os.LookupEnv(name); ok {
			fs.Set(f.Name, v)
		}
	})
}
