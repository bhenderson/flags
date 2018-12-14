// Package flags adds common extras to be used alongside the flag package
//
// Extras include
//
// * env parsing (12 factor apps) This is the only feature that sets values.
// This can be called before or after flag.Parse depending on desired
// precedence.
//
// * config parsing. The config file contains exactly what you would put on the
// command line, with support for newlines and comments. No special config
// syntax.
//
// * struct parsing for easy configs
package flags
