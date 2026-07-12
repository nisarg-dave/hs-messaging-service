package service

import "hs-messaging-service/internal/logging"

// noopLogger is a test double for logging.Logger (Liskov Substitution — SOLID L).
// It satisfies the interface so tests can inject a silent logger without
// depending on slog or producing output.
type noopLogger struct{}

func (noopLogger) Info(msg string, args ...any)  {}
func (noopLogger) Error(msg string, args ...any) {}

func testLogger() logging.Logger {
	return noopLogger{}
}
