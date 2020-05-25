package gol

import (
	"context"
	"github.com/goburrow/gol"
)

// DefaultCtxLogger implements Logger interface.
type DefaultCtxLogger struct {
	gol.Logger
}

func NewCtxLogger(name string) *DefaultCtxLogger {
	logger := gol.GetLogger(name)
	return &DefaultCtxLogger{logger}
}

// TraceCtxf logs message at Trace level.
func (logger *DefaultCtxLogger) TraceCtxf(ctx context.Context, format string, args ...interface{}) {
	logger.Tracef(format, args)
}

// DebugCtxf logs message at Debug level.
func (logger *DefaultCtxLogger) DebugCtxf(ctx context.Context, format string, args ...interface{}) {
	logger.Debugf(format, args)
}

// InfoCtxf logs message at Info level.
func (logger *DefaultCtxLogger) InfoCtxf(ctx context.Context, format string, args ...interface{}) {
	logger.Infof(format, args)
}

// WarnCtxf logs message at Warning level.
func (logger *DefaultCtxLogger) WarnCtxf(ctx context.Context, format string, args ...interface{}) {
	logger.Warnf(format, args)
}

// ErrorCtxf logs message at Error level.
func (logger *DefaultCtxLogger) ErrorCtxf(ctx context.Context, format string, args ...interface{}) {
	logger.Errorf(format, args)
}

type GAELoggerFac struct {
	projectId string
}

func (fac *GAELoggerFac) GetLogger(name string) gol.Logger {
	return NewGCLogger(fac.projectId, name)
}

func NewGAELoggerFactory(projectId string) *GAELoggerFac {
	return &GAELoggerFac{projectId}
}

type DefaultCtxLoggerFactory struct {
}

func NewDefaultCtxLoggerFactory() *DefaultCtxLoggerFactory {
	return &DefaultCtxLoggerFactory{}
}

func (fac *DefaultCtxLoggerFactory) GetLogger(name string) gol.Logger {
	return NewCtxLogger(name)
}
