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
	env := os.Environ()
	fs.VisitAll(func(f *flag.Flag) {
		name := strings.ToUpper(f.Name)
		for _, e := range env {
			if strings.HasPrefix(e, name+"=") {
				fs.Set(f.Name, e[len(name)+1:])
				break
			}
		}
	})
}
