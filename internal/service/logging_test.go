package service

// noopLogger is a test double for Logger (Liskov Substitution — SOLID L).
// It satisfies the Logger interface so tests can inject a silent logger without
// depending on slog or producing output.
type noopLogger struct{}

func (noopLogger) Info(msg string, args ...any)  {}
func (noopLogger) Error(msg string, args ...any) {}

func testLogger() Logger {
	return noopLogger{}
}
