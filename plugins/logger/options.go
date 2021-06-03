package logger

import (
	"github.com/asim/go-micro/v3/logger"
	"github.com/rs/zerolog"
)

type Mode uint8

const (
	Production Mode = iota
	Development
)

type Options struct {
	logger.Options

	// Flag for whether to log caller info (off by default)
	ReportCaller bool
	// Use this logger as system wide default logger  (off by default)
	UseAsDefault bool
	// zerolog hooks
	Hooks []zerolog.Hook
	// TimeFormat is one of time.RFC3339, time.RFC3339Nano, time.*
	TimeFormat string
	// Runtime mode. (Production by default)
	Mode Mode
	// Exit Function to call when FatalLevel log
	ExitFunc func(int)
	// Output file name
	// If set, output to a file
	OutputFilePath string
	// Output root path
	OutputRootPath string
}

type rootPathKey struct{}

func OutputRootPath(rootPath string) logger.Option {
	return logger.SetOption(rootPathKey{}, rootPath)
}

type filePathKey struct{}

func OutputFilePath(filename string) logger.Option {
	return logger.SetOption(filePathKey{}, filename)
}

type reportCallerKey struct{}

func ReportCaller() logger.Option {
	return logger.SetOption(reportCallerKey{}, true)
}

type useAsDefaultKey struct{}

func UseAsDefault() logger.Option {
	return logger.SetOption(useAsDefaultKey{}, true)
}

type developmentModeKey struct{}

func WithDevelopmentMode() logger.Option {
	return logger.SetOption(developmentModeKey{}, true)
}

type productionModeKey struct{}

func WithProductionMode() logger.Option {
	return logger.SetOption(productionModeKey{}, true)
}

type timeFormatKey struct{}

func WithTimeFormat(timeFormat string) logger.Option {
	return logger.SetOption(timeFormatKey{}, timeFormat)
}

type hooksKey struct{}

func WithHooks(hooks []zerolog.Hook) logger.Option {
	return logger.SetOption(hooksKey{}, hooks)
}

type exitKey struct{}

func WithExitFunc(exit func(int)) logger.Option {
	return logger.SetOption(exitKey{}, exit)
}
