package gol

import "github.com/goburrow/gol"
import "context"

// CtxLogger ads context related methods.
type CtxLogger interface {
	gol.Logger
	TraceCtxf(context.Context, string, ...interface{})
	DebugCtxf(context.Context, string, ...interface{})
	InfoCtxf(context.Context, string, ...interface{})
	WarnCtxf(context.Context, string, ...interface{})
	ErrorCtxf(context.Context, string, ...interface{})
}
