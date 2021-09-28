package levelled

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logs     []*log
	mu       sync.Mutex
	logger   *zap.Logger
	severity zapcore.Level
}

type log struct {
	entry  *zapcore.CheckedEntry
	fields []zap.Field
}

func NewLogger(logger *zap.Logger, severity zapcore.Level) *Logger {
	return &Logger{
		logs:     make([]*log, 0),
		logger:   logger,
		severity: severity,
	}
}

func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{
		logs:     append(make([]*log, 0), l.logs...),
		logger:   l.logger.With(fields...),
		severity: l.severity,
	}
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.handle(l.logger.Check(zapcore.DebugLevel, msg), fields...)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.handle(l.logger.Check(zapcore.InfoLevel, msg), fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.handle(l.logger.Check(zapcore.WarnLevel, msg), fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.handle(l.logger.Check(zapcore.ErrorLevel, msg), fields...)
}

func (l *Logger) DPanic(msg string, fields ...zap.Field) {
	l.handle(l.logger.Check(zapcore.DPanicLevel, msg), fields...)
}

func (l *Logger) Panic(msg string, fields ...zap.Field) {
	l.handle(l.logger.Check(zapcore.PanicLevel, msg), fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.handle(l.logger.Check(zapcore.FatalLevel, msg), fields...)
}

func (l *Logger) Erase() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.eraseLogs()
}

func (l *Logger) handle(entry *zapcore.CheckedEntry, fields ...zap.Field) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.severity <= entry.Entry.Level {
		for _, item := range l.logs {
			item.entry.Write(item.fields...)
		}

		entry.Write(fields...)

		l.eraseLogs()
	} else {
		l.logs = append(l.logs, &log{
			entry:  entry,
			fields: fields,
		})
	}
}

func (l *Logger) eraseLogs() {
	l.logs = make([]*log, 0)
}
