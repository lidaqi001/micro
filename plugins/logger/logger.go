package logger

import (
	"github.com/asim/go-micro/plugins/logger/zerolog/v3"
	"github.com/asim/go-micro/v3/logger"
	"os"
)

// Logger is a generic logging interface
//type Logger interface {
//	// Init initialises options
//	Init(options ...Option) error
//	// The Logger options
//	Options() Options
//	// Fields set fields to always be logged
//	Fields(fields map[string]interface{}) Logger
//	// Log writes a log entry
//	Log(level Level, v ...interface{})
//	// Logf writes a formatted log entry
//	Logf(level Level, format string, v ...interface{})
//	// String returns the name of logger
//	String() string
//}

var DefaultLogger logger.Logger = NewLogger("")

func Log(level logger.Level, v ...interface{}) {
	DefaultLogger.Log(level, v...)
}

func Logf(level logger.Level, format string, v ...interface{}) {
	DefaultLogger.Logf(level, format, v...)
}

func String() string {
	return DefaultLogger.String()
}

func WithFile(fileName string) (logger.Logger, error) {

	//var (
	//	f   *os.File
	//	err error
	//)
	//if !helper.IsExist(logPath) {
	//	err = helper.CreateDir(logPath)
	//	if err != nil {
	//		return nil, err
	//	}
	//}
	//
	//f, err = os.OpenFile(logPath+fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	//if err != nil {
	//	return nil, err
	//}
	//defer f.Close()
	logger.DefaultLogger = zerolog.NewLogger(logger.WithOutput(os.Stdout), zerolog.WithDevelopmentMode())
	//_ = logger.Init(zerolog.WithDevelopmentMode())
	logger.Debug(logger.DebugLevel, "test")

	return nil, nil
}
