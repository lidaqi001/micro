# 日志组件

- 基于go-micro原组件 [asim/go-micro/plugins/logger/zerolog](https://github.com/asim/go-micro/tree/master/plugins/logger/zerolog)

- 基于原logger组件接口编写，兼容原组件

- 添加配置项
    - 指定日志输出文件名
    ```text
    配置项： OutputFilePath
    参数不为空，就可以输出到对应文件（OutputFileName如果为空，则输出到os.Stdout）
    示例：logger.NewLogger(logger.OutputFilePath("debug")) 
    // 生成文件为：{OutputRootPath}/debug/20210603/15.log，格式：{OutputRootPath}/路径名/日期/小时.log
    ```

    - 指定日志根目录
    ```text
    配置项：OutputRootPath
    默认为当前目录下log目录（不存在会新建）
    示例：logger.NewLogger(logger.OutputRootPath("./log/"))
    ``` 
  
    - 配置项结构体
    ```text
    type Options struct {
    
        ······
    
        // Output file name
        // If set, output to a file
        OutputFilePath string
        // Output root path
        OutputRootPath string
    }
    ```

- 原组件需要设置Option为：**Development** 时，才会显示debug的消息，不够灵活，重写一下

```go
func (l *zeroLogger) Init(opts ...logger.Option) error {
······
	switch l.opts.Mode {
// 这里判断了环境
	case Development:
······
		//level = logger.DebugLevel
		l.zLog = zerolog.New(consOut).
// 开启了环境才输出debuglevel的日志
			Level(zerolog.DebugLevel).
			With().Timestamp().Stack().Logger()

// 默认是不输出debuglevel日志的
	default: // Production
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		l.zLog = zerolog.New(l.opts.Out).
			Level(zerolog.InfoLevel).
			With().Timestamp().Stack().Logger()
	}
······
}
```
