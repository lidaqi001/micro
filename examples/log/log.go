package main

import (
	"github.com/lidaqi001/micro/common/config"
	"github.com/lidaqi001/micro/plugins/logger"
)

func main() {
	// 默认直接输出，可以通过配置 logger.OutputFilePath 输出到文件
	// 所以除非要自定义log文件目录，否则可以直接使用，跳过设置 logger.NewLogger() 这一步
	logger.DefaultLogger = logger.NewLogger(
		// 输出日志路径
		// 设置此配置会输出到文件
		logger.OutputFilePath("test"),
		// 输出日志根目录
		// default：config.LOG_ROOT
		logger.OutputRootPath(config.LOG_ROOT),
		// 以小时分割日志
		// default：以天分割日志
		logger.SplitLogByHour(),
	)

	msg := "test log message~"
	logger.Info(msg)
	logger.Debug(msg)
	logger.Error(msg)
	logger.Log(logger.DebugLevel, msg)
}
