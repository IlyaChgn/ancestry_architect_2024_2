package utils

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"go.uber.org/zap"

	newlogger "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/server/usecases"
)

func GetFunctionName() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	values := strings.Split(frame.Function, ".")
	line := strconv.Itoa(frame.Line)
	pathParts := strings.Split(frame.File, "/")
	shortPath := strings.Join(pathParts[len(pathParts)-3:], "/")

	return values[len(values)-1] + " in " + shortPath + ":" + line
}

func GetOnlyFunctionName() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	values := strings.Split(frame.Function, ".")

	return values[len(values)-1]
}

func LogInfo(logger *zap.SugaredLogger, msg string) {
	logger = logger.With(zap.String("msg", msg))
	logger.Info()
}

func LogError(logger *zap.SugaredLogger, msg interface{}) {
	switch m := msg.(type) {
	case string:
		logger = logger.With(zap.String("msg", m))
	case error:
		logger = logger.With(zap.String("msg", m.Error()))
	default:
		logger = logger.With(zap.String("msg", "unable to convert msg to string"))
	}
	logger.Error()
}

func LogHandlerInfo(logger *zap.SugaredLogger, msg string, statusCode int) {
	logger = logger.With(zap.String("status", strconv.Itoa(statusCode)))
	logger = logger.With(zap.String("msg", msg))
	logger.Info()
}

func LogHandlerError(logger *zap.SugaredLogger, msg interface{}, statusCode int) {
	logger = logger.With(zap.String("status", strconv.Itoa(statusCode)))
	switch m := msg.(type) {
	case string:
		logger = logger.With(zap.String("msg", m))
	case error:
		logger = logger.With(zap.String("msg", m.Error()))
	default:
		logger = logger.With(zap.String("msg", "unable to convert msg to string"))
	}
	logger.Error()
}

func GetLoggerFromContext(ctx context.Context) *zap.SugaredLogger {
	loggerContextKey := fmt.Sprintf("%v", ctx.Value("loggerContextKey"))
	outputLogPath := fmt.Sprintf("%v", ctx.Value("outputLogPath"))
	errorOutputLogPath := fmt.Sprintf("%v", ctx.Value("errorOutputLogPath"))

	if logger, ok := ctx.Value(loggerContextKey).(*zap.SugaredLogger); ok {
		return logger
	}

	logger, _ := newlogger.NewLogger(strings.Split(outputLogPath, " "),
		strings.Split(errorOutputLogPath, " "))
	logger.Error("Can`t get logger from context")

	return logger
}
