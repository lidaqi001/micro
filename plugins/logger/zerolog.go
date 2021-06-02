package logger

import (
	"fmt"
	"github.com/asim/go-micro/v3/logger"
	"github.com/lidaqi001/micro/common/helper"
	"github.com/rs/zerolog"
	"os"
	"time"
)

type zlog struct {
	log zerolog.Logger
}

func (l *zlog) Log(level logger.Level, v ...interface{}) {
	l.log.WithLevel(loggerToZerologLevel(level)).Msg(fmt.Sprint(v...))
}

func (l *zlog) Logf(level logger.Level, format string, v ...interface{}) {
	l.log.WithLevel(loggerToZerologLevel(level)).Msgf(format, v...)
}

func (l *zlog) String() string {
	return "zerolog"
}

func (l *zlog) Init(options ...logger.Option) error {
	panic("implement me")
}

func (l *zlog) Options() logger.Options {
	panic("implement me")
}

func (l *zlog) Fields(fields map[string]interface{}) logger.Logger {
	panic("implement me")
}

const (
	//LOGPATH  LOGPATH/time.Now().Format(FORMAT)/*.log
	LOGPATH = "./log/"
	//FORMAT .
	FORMAT = "20060102"
)

//以天为基准,切割日志
var (
	logPath = LOGPATH + time.Now().Format(FORMAT) + "/"
)

func NewLogger(fileName string) logger.Logger {

	var (
		f   *os.File
		err error
	)
	if !helper.IsExist(logPath) {
		err = helper.CreateDir(logPath)
		if err != nil {
			return &zlog{}
		}
	}
	if len(fileName) == 0 {
		fileName = time.Now().Format(FORMAT) + ".log"
	} else {
		fileName += ".log"
	}

	f, err = os.OpenFile(logPath+fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	lg := zerolog.New(f). // f || os.Stdout
		With().Timestamp().Stack().Logger()
	z := &zlog{
		log: lg,
	}
	return z
}
