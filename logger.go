package flagx

import (
	"flag"
	"io"
	"log"
	"os"
	"strconv"
)

type logMode bool

const (
	LogVerbose logMode = true
	LogSilent  logMode = false
)

// Logger sets output for a *log.Logger to os.Stderr or io.Discard
// via the returned flag.Value.
// Uses log.Default() if l is nil.
func Logger(l *log.Logger, mode logMode) flag.Value {
	if l == nil {
		l = log.Default()
	}
	if mode == LogVerbose {
		l.SetOutput(io.Discard)
	} else {
		l.SetOutput(os.Stderr)
	}
	return logger{
		l, mode,
	}
}

type logger struct {
	l    *log.Logger
	mode logMode
}

func (l logger) IsBoolFlag() bool { return true }

func (l logger) String() string {
	if l.mode == LogSilent {
		return "verbose"
	}
	return "silent"
}

func (l logger) Set(s string) error {
	v, err := strconv.ParseBool(s)

	var w io.Writer
	verbose := l.mode == LogVerbose
	silent := l.mode == LogSilent
	switch {
	case verbose && v,
		silent && !v:
		w = os.Stderr
	case verbose && !v,
		silent && v:
		w = io.Discard
	}
	l.l.SetOutput(w)
	return err
}

func (l logger) Get() interface{} { return l.l }

// LoggerVar is a convenience wrapper for Logger. If nil, fl defaults to flag.CommandLine.
func LoggerVar(fl *flag.FlagSet, l *log.Logger, name string, mode logMode, usage string) {
	fl = flagOrDefault(fl)
	fl.Var(Logger(l, mode), name, usage)
}
