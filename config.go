package flags

import (
	"bufio"
	"flag"
	"os"
)

// Config adds a flag to f. When encountered, the named file will be split and
// given to f.Parse. Empty lines and lines beginning with # are ignored.
func Config(name, usage string, fs *flag.FlagSet) {
	if fs == nil {
		fs = flag.CommandLine
	}
	fs.Var(&config{fs: fs, name: name}, name, usage)
}

var _ flag.Value = &config{}

type config struct {
	fs     *flag.FlagSet
	name   string
	parsed bool
}

func (c *config) String() string {
	return c.name
}

func (c *config) Set(s string) error {
	// Prevent someone from puting -config="this-file"
	if c.parsed {
		return nil
	}
	file, err := os.Open(s)
	if err != nil {
		return err
	}
	defer file.Close()

	var args []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		arg := scanner.Text()
		if arg == "" || arg[0] == '#' {
			continue
		}
		args = append(args, arg)
	}
	if err = scanner.Err(); err != nil {
		return err
	}
	// continue where we left off
	defer c.fs.Parse(c.fs.Args())
	defer c.setParsed()()
	return c.fs.Parse(args)
}

func (c *config) setParsed() func() {
	c.parsed = true
	return func() { c.parsed = false }
}
