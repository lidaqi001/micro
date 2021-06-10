package logger

import (
	"context"
	"github.com/rs/zerolog"
	"io"
)

type Mode uint8

const (
	Production Mode = iota
	Development
)

type Option func(*Options)

type Options struct {
	// It's common to set this to a file, or leave it default which is `os.Stderr`
	Out io.Writer
	// The logging level the logger should log at. default is `InfoLevel`
	Level Level
	// fields to always be logged
	Fields  map[string]interface{}
	Context context.Context

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
	// Split the log by the hour
	SplitLogByHour bool
}

type splitLogByHourKey struct{}

func SplitLogByHour() Option {
	return SetOption(splitLogByHourKey{}, true)
}

type rootPathKey struct{}

func OutputRootPath(rootPath string) Option {
	return SetOption(rootPathKey{}, rootPath)
}

type filePathKey struct{}

func OutputFilePath(filename string) Option {
	return SetOption(filePathKey{}, filename)
}

type reportCallerKey struct{}

func ReportCaller() Option {
	return SetOption(reportCallerKey{}, true)
}

type useAsDefaultKey struct{}

func UseAsDefault() Option {
	return SetOption(useAsDefaultKey{}, true)
}

type developmentModeKey struct{}

func WithDevelopmentMode() Option {
	return SetOption(developmentModeKey{}, true)
}

type productionModeKey struct{}

func WithProductionMode() Option {
	return SetOption(productionModeKey{}, true)
}

type timeFormatKey struct{}

func WithTimeFormat(timeFormat string) Option {
	return SetOption(timeFormatKey{}, timeFormat)
}

type hooksKey struct{}

func WithHooks(hooks []zerolog.Hook) Option {
	return SetOption(hooksKey{}, hooks)
}

type exitKey struct{}

func WithExitFunc(exit func(int)) Option {
	return SetOption(exitKey{}, exit)
}
