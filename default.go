package gol

import (
	"context"
	"github.com/goburrow/gol"
)

// DefaultCtxLogger implements Logger interface.
type DefaultCtxLogger struct {
	gol.DefaultLogger
}

func NewCtxLogger(name string, parent *DefaultCtxLogger) *DefaultCtxLogger {
	logger := gol.New(name, nil)
	return &DefaultCtxLogger{*logger}
}

// TraceCtxf logs message at Trace level.
func (logger *DefaultCtxLogger) TraceCtxf(ctx context.Context, format string, args ...interface{}) {
	logger.Printf(gol.Trace, format, args)
}

// DebugCtxf logs message at Debug level.
func (logger *DefaultCtxLogger) DebugCtxf(ctx context.Context, format string, args ...interface{}) {
	logger.Printf(gol.Debug, format, args)
}

// InfoCtxf logs message at Info level.
func (logger *DefaultCtxLogger) InfoCtxf(ctx context.Context, format string, args ...interface{}) {
	logger.Printf(gol.Info, format, args)
}

// WarnCtxf logs message at Warning level.
func (logger *DefaultCtxLogger) WarnCtxf(ctx context.Context, format string, args ...interface{}) {
	logger.Printf(gol.Warn, format, args)
}

// ErrorCtxf logs message at Error level.
func (logger *DefaultCtxLogger) ErrorCtxf(ctx context.Context, format string, args ...interface{}) {
	logger.Printf(gol.Error, format, args)
}

type DefaultCtxLoggerFactory struct {
	projectId string
}

func NewDefaultCtxLoggerFactory(projectId string) *DefaultCtxLoggerFactory {
	return &DefaultCtxLoggerFactory{projectId}
}

func (fac *DefaultCtxLoggerFactory) GetLogger(name string) gol.Logger {
	return NewGCLogger(fac.projectId, name)
}
