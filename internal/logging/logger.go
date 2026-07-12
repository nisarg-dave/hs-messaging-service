package logging

// Logger is the shared logging contract for cross-cutting components.
//
// Design: Dependency Inversion (SOLID — D) — services and middleware depend on
// this abstraction, not on *slog.Logger or a global singleton. Production
// passes a concrete *slog.Logger from main; tests pass a no-op fake.
//
// Design: Interface Segregation (SOLID — I) — only the methods callers actually
// use are on the interface, not the full slog API.
//
// This package is intentionally neutral (not owned by service or api) so both
// layers can import it without inverting dependencies.
type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}
