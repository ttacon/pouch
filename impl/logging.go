package impl

import (
	"log"
	"os"
	"time"

	"github.com/ttacon/chalk"
)

type Logger interface {
	Print(v ...interface{})
}

type logger struct {
	l *log.Logger
}

var defaultLogger = logger{
	l: log.New(os.Stdout, "", 0),
}

func (l logger) Print(v ...interface{}) {
	if len(v) > 1 {
		currentTime := chalk.Yellow.Color(time.Now().Format("2006-02-01 15:04:05"))
		v = append([]interface{}{currentTime}, v...)
		l.l.Println(v...)
	}
}
