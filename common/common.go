package common

import (
	"github.com/asim/go-micro/v3/logger"
	"github.com/lidaqi001/micro/common/config"
	log "github.com/lidaqi001/micro/plugins/logger"
)

/*******************************************************
			            设置配置
*******************************************************/

// 设置默认日志记录器
func SetDefaultLoggerForZerolog(defaultPath string) {
	logger.DefaultLogger = log.NewLogger(
		// 日志目录
		log.OutputFilePath(defaultPath),
		// 日志根目录
		log.OutputRootPath(config.LOG_ROOT),
	)
}
