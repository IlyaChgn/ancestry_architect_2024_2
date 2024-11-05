package utils

import (
	"context"
	"runtime"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/config"
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
	if logger, ok := ctx.Value(config.LoggerContextKey).(*zap.SugaredLogger); ok {
		return logger
	}

	logger, _ := newlogger.NewLogger(strings.Split(config.OutputLogPath, " "),
		strings.Split(config.ErrorOutputLogPath, " "))
	logger.Error("Couldnt get logger from context")

	return logger
}
