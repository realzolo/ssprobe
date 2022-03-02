package util

import (
	"log"
)

type Logger struct {
}

func (l *Logger) OnlyLog(message string) {
	log.Println(message)
}

func (l *Logger) LogWithError(err error, message string) {
	if err != nil {
		log.Println(message, err.Error())
	}
}

func (l *Logger) LogWithFormat(format string, v ...interface{}) {
	log.Printf(format, v)
}

func (l *Logger) ErrorWithFormat(err error, format string, v ...interface{}) {
	if err != nil {
		log.Printf(format, v)
	}
}

func (l *Logger) LogWithExit(message string) {
	log.Fatalln(message)
}

func (l *Logger) ErrorWithExit(err error, message string) {
	if err != nil {
		log.Fatalln(message, err.Error())
	}
}
