package logger

import (
	"context"
	"fmt"
	"github.com/asim/go-micro/v3/logger"
	"github.com/lidaqi001/micro/common/helper"
	"github.com/rs/zerolog"
	zelog "github.com/rs/zerolog/log"
	"os"
	"time"
)

type zlog struct {
	log  zerolog.Logger
	opts Options
}

func (l *zlog) Log(level logger.Level, v ...interface{}) {
	l.log.WithLevel(loggerToZerologLevel(level)).Msg(fmt.Sprint(v...))
	// Invoke os.Exit because unlike zerolog.Logger.Fatal zerolog.Logger.WithLevel won't stop the execution.
	if level == logger.FatalLevel {
		l.opts.ExitFunc(1)
	}
}

func (l *zlog) Logf(level logger.Level, format string, v ...interface{}) {
	l.log.WithLevel(loggerToZerologLevel(level)).Msgf(format, v...)
	// Invoke os.Exit because unlike zerolog.Logger.Fatal zerolog.Logger.WithLevel won't stop the execution.
	if level == logger.FatalLevel {
		l.opts.ExitFunc(1)
	}
}

func (l *zlog) String() string {
	return "zerolog"
}

func (l *zlog) Init(opts ...logger.Option) error {
	for _, o := range opts {
		o(&l.opts.Options)
	}

	if hs, ok := l.opts.Context.Value(hooksKey{}).([]zerolog.Hook); ok {
		l.opts.Hooks = hs
	}
	if tf, ok := l.opts.Context.Value(timeFormatKey{}).(string); ok {
		l.opts.TimeFormat = tf
	}
	if exitFunction, ok := l.opts.Context.Value(exitKey{}).(func(int)); ok {
		l.opts.ExitFunc = exitFunction
	}
	if caller, ok := l.opts.Context.Value(reportCallerKey{}).(bool); ok && caller {
		l.opts.ReportCaller = caller
	}
	if useDefault, ok := l.opts.Context.Value(useAsDefaultKey{}).(bool); ok && useDefault {
		l.opts.UseAsDefault = useDefault
	}
	if devMode, ok := l.opts.Context.Value(developmentModeKey{}).(bool); ok && devMode {
		l.opts.Mode = Development
	}
	if prodMode, ok := l.opts.Context.Value(productionModeKey{}).(bool); ok && prodMode {
		l.opts.Mode = Production
	}

	if filePath, ok := l.opts.Context.Value(filePathKey{}).(string); ok {
		l.opts.OutputFilePath = filePath
	}
	if rootPath, ok := l.opts.Context.Value(rootPathKey{}).(string); ok {
		l.opts.OutputRootPath = rootPath
	}
	if byHour, ok := l.opts.Context.Value(splitLogByHourKey{}).(bool); ok && byHour {
		l.opts.SplitLogByHour = byHour
	}

	// RESET
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.ErrorStackMarshaler = nil
	zerolog.CallerSkipFrameCount = 4

	//switch l.opts.Mode {
	//case Development:
	//default: // Production
	//}

	// Set Log
	var err error
	if l.log, err = getZlog(l.opts); err != nil {
		return err
	}

	// Set log Level if not default
	if l.opts.Level != 100 {
		zerolog.SetGlobalLevel(loggerToZerologLevel(l.opts.Level))
		l.log = l.log.Level(loggerToZerologLevel(l.opts.Level))
	}

	// Reporting caller
	if l.opts.ReportCaller {
		l.log = l.log.With().Caller().Logger()
	}

	// Adding hooks if exist
	for _, hook := range l.opts.Hooks {
		l.log = l.log.Hook(hook)
	}

	// Setting timeFormat
	if len(l.opts.TimeFormat) > 0 {
		zerolog.TimeFieldFormat = l.opts.TimeFormat
	}

	// Adding seed fields if exist
	if l.opts.Fields != nil {
		l.log = l.log.With().Fields(l.opts.Fields).Logger()
	}

	// Also set it as zerolog's Default logger
	if l.opts.UseAsDefault {
		zelog.Logger = l.log
	}

	return nil
}

func (l *zlog) Options() logger.Options {
	return l.opts.Options
}

func (l *zlog) Fields(fields map[string]interface{}) logger.Logger {
	panic("implement me")
}

// NewLogger builds a new logger based on options
func NewLogger(opts ...logger.Option) logger.Logger {
	// Default options
	options := Options{
		Options: logger.Options{
			Level:   100,
			Fields:  make(map[string]interface{}),
			Out:     os.Stderr,
			Context: context.Background(),
		},
		ReportCaller:   false,
		UseAsDefault:   false,
		Mode:           Production,
		ExitFunc:       os.Exit,
		OutputFilePath: "",
		SplitLogByHour: false,
		SplitLogByDay:  false,
	}

	l := &zlog{opts: options}
	_ = l.Init(opts...)
	return l
}

//FORMAT Date
const FORMAT = "20060102"

// Output to the root path
var ORootPath = "./log/"

func getZlog(opts Options) (zerolog.Logger, error) {

	var log zerolog.Logger

	filePath := opts.OutputFilePath

	// If the filePath is not set , default output to os.Stdout
	if len(filePath) == 0 {
		log = zerolog.New(os.Stdout).With().Timestamp().Stack().Logger()
		return log, nil
	}

	if len(opts.OutputRootPath) > 0 {
		ORootPath = opts.OutputRootPath
	}

	var (
		err  error
		file *os.File

		date    = time.Now().Format(FORMAT)
		logPath = ORootPath + filePath
	)

	var path string
	if opts.SplitLogByHour {
		// 以小时分割日志
		// 以日期作为文件夹名，小时作为文件名
		logPath += "/" + date
		path = fmt.Sprintf("%s/%d.log", logPath, time.Now().Hour())
	} else {
		// (默认) 以日期分割日志
		path = fmt.Sprintf("%s/%s.log", logPath, date)
	}

	// 检测路径，不存在就新增
	if !helper.IsExist(logPath) {
		err = helper.CreateDir(logPath)
		if err != nil {
			return zerolog.Logger{}, err
		}
	}

	file, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	// f || os.Stdout
	log = zerolog.New(file).With().Timestamp().Stack().Logger()
	return log, nil
}
