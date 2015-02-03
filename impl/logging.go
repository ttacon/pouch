package impl

import (
	"log"
	"os"
	"time"

	"github.com/ttacon/chalk"
)

type Logger interface {
	Print(v ...interface{})
	V(level VerbosityLevel) Logger
	SetVerbosity(level VerbosityLevel)
}

type VerbosityLevel int

const (
	TRACE VerbosityLevel = iota
	INFO
	WARN
	ERROR
)

type logger struct {
	l             *log.Logger
	currVerbosity VerbosityLevel
}

var defaultVerbosity = INFO

func SetDefaultVerbosity(level VerbosityLevel) {
	defaultVerbosity = level
}

func defaultLogger() Logger {
	return &logger{
		l:             log.New(os.Stdout, "", 0),
		currVerbosity: defaultVerbosity,
	}
}

func (l *logger) Print(v ...interface{}) {
	currentTime := chalk.Yellow.Color(time.Now().Format("2006-02-01 15:04:05"))
	v = append([]interface{}{currentTime}, v...)
	l.l.Println(v...)
}

func (l *logger) V(level VerbosityLevel) Logger {
	return verbosityLogger(l.currVerbosity == level)
}

func (l *logger) SetVerbosity(level VerbosityLevel) { l.currVerbosity = level }

type verbosityLogger bool

func (vl verbosityLogger) Print(v ...interface{}) {
	if vl {
		currentTime := chalk.Yellow.Color(time.Now().Format("2006-02-01 15:04:05"))
		v = append([]interface{}{currentTime}, v...)
		// TODO(ttacon): don't be a baby and not use the Logger
		// this is from, this will need to change when we allow custom Loggers
		log.Println(v...)
	}
}

func (v verbosityLogger) V(level VerbosityLevel) Logger     { return v }
func (v verbosityLogger) SetVerbosity(level VerbosityLevel) {}
