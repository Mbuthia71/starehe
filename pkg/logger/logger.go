package logger

import (
	"log"
	"os"
)

type Logger struct {
	*log.Logger
	info  *log.Logger
	error *log.Logger
	debug *log.Logger
	warn  *log.Logger
}

func New() *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lshortfile),
		info:   log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lshortfile),
		error:  log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lshortfile),
		debug:  log.New(os.Stdout, "[DEBUG] ", log.LstdFlags|log.Lshortfile),
		warn:   log.New(os.Stdout, "[WARN] ", log.LstdFlags|log.Lshortfile),
	}
}

func (l *Logger) Info(v ...interface{}) {
	l.info.Println(v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.error.Println(v...)
}

func (l *Logger) Debug(v ...interface{}) {
	l.debug.Println(v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.info.Printf(format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.error.Printf(format, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.debug.Printf(format, v...)
}

func (l *Logger) Warn(v ...interface{}) {
	l.warn.Println(v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.warn.Printf(format, v...)
}
