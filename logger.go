package runway

import (
	"context"
	"log/slog"
	"os"
	"time"

	gormLog "gorm.io/gorm/logger"
)

var logger *RunwayLogger

// RunwayLogger is a composite logger that embeds a slog.Logger for structured logging
// and includes a gormLogger for database logging integration.
// This allows unified logging across application and database layers.
type RunwayLogger struct {
	*slog.Logger
	dbLogger   *slog.Logger
	gormLogger *gormLogger
}

func DisableDatabaseLogs() {
	logger.gormLogger.Level = gormLog.Silent
}

func NewRunwayLogger() *RunwayLogger {
	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	db := l.With("module", "database")

	r := &RunwayLogger{
		Logger:     l,
		dbLogger:   db,
		gormLogger: &gormLogger{Level: gormLog.Info},
	}

	return r
}

// gormLogger uses slog under the hood and imllements the gorm logger interface.
// So it can be reused for all GORM outputs as well.
//
// More details can be found here: https://gorm.io/docs/logger.html
//
//	type Interface interface {
//		LogMode(LogLevel) Interface
//		Info(context.Context, string, ...interface{})
//		Warn(context.Context, string, ...interface{})
//		Error(context.Context, string, ...interface{})
//		Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error)
//	}
type gormLogger struct {
	Level gormLog.LogLevel
}

func (l gormLogger) LogMode(level gormLog.LogLevel) gormLog.Interface {
	l.Level = level
	return l
}

func (l gormLogger) Info(ctx context.Context, message string, details ...interface{}) {
	logger.Info(message, slog.Any("details", detailsToSlog(details)))
}

func (l gormLogger) Warn(ctx context.Context, message string, details ...interface{}) {
	logger.Warn(message, slog.Any("details", detailsToSlog(details)))
}

func (l gormLogger) Error(ctx context.Context, message string, details ...interface{}) {
	logger.Error(message, slog.Any("details", detailsToSlog(details)))
}

func (l gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rowsAffected := fc()

	if err != nil {
		logger.Error("sql operation", "begin", begin, "sql", sql, "rows_affected", rowsAffected, "error", err)
	} else {
		logger.Error("sql operation", "begin", begin, "sql", sql, "rows_affected", rowsAffected)
	}
}

func detailsToSlog(details ...any) []slog.Attr {
	list := make([]slog.Attr, 0, len(details))

	for i, d := range details {
		list = append(list, slog.Any("detail_"+string(i-1), d))
	}

	return list
}
