package namespace

import (
	"context"
	"log/slog"
)

type Logger interface {
	Println(ctx context.Context, msg string, args ...any)
}

type StdLogger struct{}

func (StdLogger) Println(_ context.Context, msg string, args ...any) {
	slog.Info(msg, args...)
}
