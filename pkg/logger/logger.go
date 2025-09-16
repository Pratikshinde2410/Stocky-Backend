package logger

import (
    "go.uber.org/zap"
)

type Logger interface {
    Info(msg string, keysAndValues ...interface{})
    Error(msg string, keysAndValues ...interface{})
    Debug(msg string, keysAndValues ...interface{})
}

type zapLogger struct {
    l *zap.SugaredLogger
}

func New() Logger {
    base, _ := zap.NewProduction()
    return &zapLogger{l: base.Sugar()}
}

func (z *zapLogger) Info(msg string, keysAndValues ...interface{})  { z.l.Infow(msg, keysAndValues...) }
func (z *zapLogger) Error(msg string, keysAndValues ...interface{}) { z.l.Errorw(msg, keysAndValues...) }
func (z *zapLogger) Debug(msg string, keysAndValues ...interface{}) { z.l.Debugw(msg, keysAndValues...) }



