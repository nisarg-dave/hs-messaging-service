package service

// Shared fixture UUIDs for service package tests.
// Files ending in _test.go in the same package can call these freely
// (message_service_test.go and conversation_service_test.go both use them).
func validUUID() string { return "11111111-1111-1111-1111-111111111111" }
func otherUUID() string { return "22222222-2222-2222-2222-222222222222" }
