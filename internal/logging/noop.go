package logging

// noopLogger is a silent Logger implementation used in tests.
//
// Design: Liskov Substitution (SOLID — L) — any code that depends on Logger can
// accept this in place of *slog.Logger without changing call sites. Methods are
// empty so unit tests stay quiet and don't depend on slog output.
type noopLogger struct{}

func (noopLogger) Info(msg string, args ...any)  {}
func (noopLogger) Error(msg string, args ...any) {}

// Noop returns a Logger that discards all messages.
// Lives in this package (not a _test.go file) so other packages' tests can
// import it — Go test helpers in _test.go files are not importable.
func Noop() Logger {
	return noopLogger{}
}
