package gol

import (
	"cloud.google.com/go/logging"
	"context"
	"encoding/json"
	"fmt"
	"github.com/goburrow/gol"
	"log"
	"net/http"
	"strings"
)

const cloudTraceContext = "X-Cloud-Trace-Context"

func init() {
	// Disable log prefixes such as the default timestamp.
	// Prefix text prevents the message from being parsed as JSON.
	// A timestamp is added when shipping logs to Stackdriver.
	log.SetFlags(0)
}

// Entry defines a log entry.
type entry struct {
	Message  string `json:"message"`
	Severity string `json:"severity,omitempty"`

	// Stackdriver Log Viewer allows filtering and display of this as `jsonPayload.component`.
	Component string `json:"component,omitempty"`
}

// String renders an entry structure to the JSON format expected by Stackdriver.
func (e entry) String() string {
	if e.Severity == "" {
		e.Severity = "INFO"
	}
	out, err := json.Marshal(e)
	if err != nil {
		log.Printf("json.Marshal: %v", err)
	}
	return string(out)
}

// Google Cloud Logger see: https://cloud.google.com/logging
type GCLogger struct {
	logName       string
	projectID     string
	level         gol.Level
	componentName string
}

func NewGCLogger(projectID, logName string, componentName string) *GCLogger {
	return &GCLogger{
		logName:       logName,
		componentName: componentName,
		projectID:     projectID,
		level:         gol.Uninitialized,
	}
}

// Tracef logs message at Trace level.
func (logger *GCLogger) Tracef(format string, args ...interface{}) {
	logger.Printf(gol.Trace, format, args)
}

// TraceCtxf logs message at Trace level.
func (logger *GCLogger) TraceCtxf(ctx context.Context, format string, args ...interface{}) {
	logger.PrintCtxf(ctx, gol.Trace, format, args)
}

// TraceEnabled checks if Trace level is enabled.
func (logger *GCLogger) TraceEnabled() bool {
	return logger.loggable(gol.Trace)
}

// Debugf logs message at Debug level.
func (logger *GCLogger) Debugf(format string, args ...interface{}) {
	logger.Printf(gol.Debug, format, args)
}

// DebugCtxf logs message at Trace level.
func (logger *GCLogger) DebugCtxf(ctx context.Context, format string, args ...interface{}) {
	logger.PrintCtxf(ctx, gol.Debug, format, args)
}

// DebugEnabled checks if Debug level is enabled.
func (logger *GCLogger) DebugEnabled() bool {
	return logger.loggable(gol.Debug)
}

// Infof logs message at Info level.
func (logger *GCLogger) Infof(format string, args ...interface{}) {
	logger.Printf(gol.Info, format, args)
}

// InfoCtxf logs message at Trace level.
func (logger *GCLogger) InfoCtxf(ctx context.Context, format string, args ...interface{}) {
	logger.PrintCtxf(ctx, gol.Info, format, args)
}

// InfoEnabled checks if Info level is enabled.
func (logger *GCLogger) InfoEnabled() bool {
	return logger.loggable(gol.Info)
}

// Warnf logs message at Warning level.
func (logger *GCLogger) Warnf(format string, args ...interface{}) {
	logger.Printf(gol.Warn, format, args)
}

// WarnCtxf logs message at Trace level.
func (logger *GCLogger) WarnCtxf(ctx context.Context, format string, args ...interface{}) {
	logger.PrintCtxf(ctx, gol.Warn, format, args)
}

// WarnEnabled checks if Warning level is enabled.
func (logger *GCLogger) WarnEnabled() bool {
	return logger.loggable(gol.Warn)
}

// Errorf logs message at Error level.
func (logger *GCLogger) Errorf(format string, args ...interface{}) {
	logger.Printf(gol.Error, format, args)
}

// ErrorCtxf logs message at Trace level.
func (logger *GCLogger) ErrorCtxf(ctx context.Context, format string, args ...interface{}) {
	logger.PrintCtxf(ctx, gol.Error, format, args)
}

// ErrorEnabled checks if Error level is enabled.
func (logger *GCLogger) ErrorEnabled() bool {
	return logger.loggable(gol.Error)
}

// Level returns level of this logger or parent if not set.
func (logger *GCLogger) Level() gol.Level {
	for logger != nil {
		if logger.level != gol.Uninitialized {
			return logger.level
		}
	}
	return gol.Off
}

// SetLevel changes logging level of this logger.
func (logger *GCLogger) SetLevel(level gol.Level) {
	logger.level = level
}

// loggable checks if the given logging level is enabled within this logger.
func (logger *GCLogger) loggable(level gol.Level) bool {
	return level >= logger.Level()
}

// log performs logging with given parameters.
func (logger *GCLogger) Printf(level gol.Level, format string, args []interface{}) {
	logger.PrintCtxf(context.TODO(), level, format, args)
}

var gclLevelStrings = map[gol.Level]string{
	gol.Trace: "DEFAULT",
	gol.Debug: "DEBUG",
	gol.Info:  "INFO",
	gol.Warn:  "WARNING",
	gol.Error: "ERROR",
}

func (logger *GCLogger) PrintCtxf(ctx context.Context, trace gol.Level, format string, args []interface{}) {
	if trace < logger.level {
		return
	}
	entry := entry{
		Message:   fmt.Sprintf(format, args),
		Component: logger.componentName,
	}

	gcEntry := logging.Entry{
		Severity: logging.ParseSeverity(gclLevelStrings[trace]),
		Payload:  entry}

	if tokenStr, ok := ctx.Value(cloudTraceContext).(string); ok {
		gcEntry.Trace = tokenStr
	}

	// Creates a client.
	client, err := logging.NewClient(ctx, logger.projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	gcLogger := client.Logger(logger.logName)
	gcLogger.Log(gcEntry)
}

// Handler to add request context to logs see: https://cloud.google.com/endpoints/docs/openapi/tracing
func (logger *GCLogger) GCLoggerHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var trace string
		if logger.projectID != "" {
			traceHeader := r.Header.Get("X-Cloud-Trace-Context")
			traceParts := strings.Split(traceHeader, "/")
			if len(traceParts) > 0 && len(traceParts[0]) > 0 {
				trace = fmt.Sprintf("projects/%s/traces/%s", logger.projectID, traceParts[0])
			}
		}
		if trace == "" {
			h.ServeHTTP(w, r)
		} else {
			ctx := r.Context()
			ctx = context.WithValue(ctx, cloudTraceContext, trace)
			h.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
