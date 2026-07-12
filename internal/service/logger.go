package service

// Logger is the logging behavior services need.
//
// Design: Dependency Inversion (SOLID — D) — services depend on this
// abstraction, not on *slog.Logger or a global singleton. Production passes
// a concrete *slog.Logger from main; tests pass a no-op fake.
//
// Design: Interface Segregation (SOLID — I) — only the two methods services
// actually call are on the interface, not the full slog API.
//
// Pattern: same consumer-side interface approach as MessageRepository — the
// type that uses the dependency defines the contract.
type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}
