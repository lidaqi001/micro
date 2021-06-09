package main

import (
	"github.com/asim/go-micro/v3/logger"
	"github.com/lidaqi001/micro/common/config"
	log "github.com/lidaqi001/micro/plugins/logger"
)

func main() {
	// 会在服务初始化时，初始化log目录
	// 所以除非要自定义log文件目录，否则可以直接使用，跳过设置 log.NewLogger() 这一步
	logger.DefaultLogger = log.NewLogger(
		// 日志目录
		log.OutputFilePath("test"),
		// 日志根目录
		log.OutputRootPath(config.LOG_ROOT),
		// 以小时分割日志
		log.SplitLogByHour(),
	)

	msg := "test log message~"
	logger.Info(msg)
	logger.Debug(msg)
	logger.Error(msg)
	logger.Log(logger.DebugLevel, msg)
}
