package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Logger struct {
	prefix string
}

const (
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorCyan    = "\033[36m"
	ColorWhite   = "\033[37m"
	ColorMagenta = "\033[35m"
)

func NewLogger(prefix string) *Logger {
	log.SetOutput(os.Stdout)
	log.SetFlags(0) // remove data/hora padrão
	return &Logger{prefix: prefix}
}

func (l *Logger) log(level string, color string, msg string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	formattedMsg := fmt.Sprintf(msg, args...)
	logLine := fmt.Sprintf("%s%-19s [%s] %-5s %s%s\n", color, timestamp, l.prefix, level, formattedMsg, ColorReset)

	if level == "FATAL" {
		log.Fatal(logLine)
	} else {
		log.Print(logLine)
	}
}

func (l *Logger) Start(msg string, args ...interface{})   { l.log("START", ColorCyan, msg, args...) }
func (l *Logger) Info(msg string, args ...interface{})    { l.log("INFO", ColorWhite, msg, args...) }
func (l *Logger) Warn(msg string, args ...interface{})    { l.log("WARN", ColorYellow, msg, args...) }
func (l *Logger) Error(msg string, args ...interface{})   { l.log("ERROR", ColorRed, msg, args...) }
func (l *Logger) Success(msg string, args ...interface{}) { l.log("SUCCESS", ColorGreen, msg, args...) }
func (l *Logger) Fatal(msg string, args ...interface{})   { l.log("FATAL", ColorMagenta, msg, args...) }
