package logger

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
	"strings"
)

var (
	logger *logrus.Entry
)

// Initialize creates a singleton logrus.Entry based on the parameters.
// It reads the env to set up the formatter and service to add the fields in all log entries.
func Initialize(config, service string) *logrus.Entry {
	log := logrus.New()

	// Customize the JSON formatter to use `message` instead of default `msg`
	log.Formatter = &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg:  "message",
			logrus.FieldKeyTime: "time",
		},
	}

	// Default environment - development
	env := "development"
	// Check if config contains "stg" and set env to "staging" accordingly
	if strings.Contains(config, "stg") {
		env = "staging"
	} else if strings.Contains(config, "prd") {
		env = "production"
	}
	config = strings.Trim(config, ".json")
	logger = log.WithFields(logrus.Fields{
		"env":     env,
		"config":  config,
		"service": service,
	})
	return logger
}

func GetLogger() *logrus.Entry {
	return logger
}

func GetLoggerWithContext(ctx context.Context) *logrus.Entry {

	// Extract trace context from the incoming context
	otelSpan := trace.SpanFromContext(ctx)
	otelSpanContext := otelSpan.SpanContext()

	logRefId := ctx.Value("x-request-id")
	//will get from request log
	if logRefId == nil {
		// will get from internal processes
		logRefId = ctx.Value("x-process-id")
	}
	if logRefId == nil {
		logRefId = uuid.New()
	}
	return logger.WithFields(logrus.Fields{
		"req-id":    logRefId,
		"traceInfo": otelSpanContext,
	})
}

func GetLogInstance(ctx context.Context, opts ...string) *logrus.Entry {
	// Extract trace context from the incoming context
	otelSpan := trace.SpanFromContext(ctx)
	otelSpanContext := otelSpan.SpanContext()
	var functionName string
	if len(opts) > 0 {
		functionName = opts[0]
	}
	// Include trace information and function name in the log entry
	logEntry := logger.WithFields(logrus.Fields{
		"traceInfo": otelSpanContext,
		"function":  functionName,
	})
	return logEntry
}
