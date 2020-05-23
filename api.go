package gol

import "context"

// Logger specifies how logging in application is done.
type Logger interface {
	Tracef(string, ...interface{})
	TraceEnabled() bool
	Debugf(string, ...interface{})
	DebugEnabled() bool
	Infof(string, ...interface{})
	InfoEnabled() bool
	Warnf(string, ...interface{})
	WarnEnabled() bool
	Errorf(string, ...interface{})
	ErrorEnabled() bool
}

// CtxLogger ads context related methods.
type CtxLogger interface {
	Logger
	TraceCtxf(context.Context, string, ...interface{})
	DebugCtxf(context.Context, string, ...interface{})
	InfoCtxf(context.Context, string, ...interface{})
	WarnCtxf(context.Context, string, ...interface{})
	ErrorCtxf(context.Context, string, ...interface{})
}

// Factory produces Logger.
type Factory interface {
	GetLogger(name string) Logger
}
