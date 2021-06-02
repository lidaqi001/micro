# 日志组件

- 基于go-micro原组件 [asim/go-micro/plugins/logger/zerolog](https://github.com/asim/go-micro/tree/master/plugins/logger/zerolog)

- 基于原logger组件接口编写，兼容原组件

- 原组件设置需要环境变量为：**Development** 时，才会显示debug的消息，不够灵活，重写一下


```go
func (l *zeroLogger) Init(opts ...logger.Option) error {
	
······

	switch l.opts.Mode {
// 这里判断了环境
	case Development:
		zerolog.ErrorStackMarshaler = func(err error) interface{} {
			fmt.Println(string(debug.Stack()))
			return nil
		}
		consOut := zerolog.NewConsoleWriter(
			func(w *zerolog.ConsoleWriter) {
				if len(l.opts.TimeFormat) > 0 {
					w.TimeFormat = l.opts.TimeFormat
				}
				w.Out = l.opts.Out
				w.NoColor = false
			},
		)
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