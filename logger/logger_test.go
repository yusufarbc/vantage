package logger

import (
	"context"
	"log/slog"
	"testing"
)

func TestLogLevel(t *testing.T) {
	tests := map[string]slog.Level{
		"":      slog.LevelInfo,
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}
	config := &Config{}
	ctx := context.Background()
	for level, expected := range tests {
		config.Level = level
		err := Setup(config)
		if err != nil {
			t.Fatalf("error setting logging level %v", err)
		}
		if !Logger.Enabled(ctx, expected) {
			t.Fatalf("invalid logging level for config %q: expected %v to be enabled", level, expected)
		}
		if expected > slog.LevelDebug && Logger.Enabled(ctx, expected-1) {
			t.Fatalf("invalid logging level for config %q: expected %v to be disabled", level, expected-1)
		}
	}
}
