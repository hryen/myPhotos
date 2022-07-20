package logger

import (
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func init() {
	// 日志
	flags := log.Ldate | log.Lmicroseconds | log.Lshortfile
	InfoLogger = log.New(os.Stdout, "INFO  ", flags)
	ErrorLogger = log.New(os.Stderr, "ERROR ", flags)
}
