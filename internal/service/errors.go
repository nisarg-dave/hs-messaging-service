package service

import (
	"errors"
	"fmt"
)

// maxContentLength caps the byte length of a message Content field.
const maxContentLength = 4000

// ErrValidation is the umbrella sentinel for validation failures returned by
// the service layer. Every specific validation error below wraps it using %w,
// and each service method wraps the result one more time with an operation
// prefix ("create message:", "list conversations:", ...).
//
// The chain ends up nested like an onion:
//
//	ErrValidation                          (sentinel,  innermost)
//	  -> errInvalidUUID                    ("validation error: must be a valid UUID")
//	    -> "create message: ..."           (added by MessageService.CreateMessage)
//
// Flow across layers for, say, a bad SenderID UUID:
//
//  1. Service returns:
//     fmt.Errorf("create message: %w", errInvalidUUID)
//     -> err.Error() = "create message: validation error: must be a valid UUID"
//
//  2. Handler runs writeServiceError(c, err). Because every link in the chain
//     was wrapped with %w, errors.Is(err, service.ErrValidation) walks the
//     chain, finds ErrValidation at the bottom, and returns true.
//
//  3. Handler responds with HTTP 400 and JSON body:
//     {"error": "create message: validation error: must be a valid UUID"}
//
// Non-validation errors (e.g. a raw GORM/DB error from the repo) don't wrap
// ErrValidation, so errors.Is returns false and the handler falls back to 500.
var ErrValidation = errors.New("validation error")

// ErrNotFound is the sentinel for "you asked for a resource that doesn't
// exist". The service translates repository-specific not-found errors (e.g.
// gorm.ErrRecordNotFound) into this so handlers can map it to HTTP 404
// without importing GORM (forbidden by the layered-architecture rule).
var ErrNotFound = errors.New("not found")

var errMessageNotFound = fmt.Errorf("%w: message", ErrNotFound)

var (
	errEmptyUserID      = fmt.Errorf("%w: userID is required", ErrValidation)
	errEmptyOtherID     = fmt.Errorf("%w: other userID is required", ErrValidation)
	errEmptySenderID    = fmt.Errorf("%w: senderId is required", ErrValidation)
	errEmptyRecipientID = fmt.Errorf("%w: recipientId is required", ErrValidation)
	errEmptyContent     = fmt.Errorf("%w: content is required", ErrValidation)
	errEmptyMessageID   = fmt.Errorf("%w: messageId is required", ErrValidation)
	errInvalidUUID      = fmt.Errorf("%w: must be a valid UUID", ErrValidation)
	errSelfMessage      = fmt.Errorf("%w: senderId and recipientId must differ", ErrValidation)
	errSelfConversation = fmt.Errorf("%w: cannot fetch conversation with self", ErrValidation)
	errContentTooLong   = fmt.Errorf("%w: content exceeds %d bytes", ErrValidation, maxContentLength)
)
