// Package kb is the domain core of the knowledger knowledge-base engine:
// recursive type-less content nodes, dual-id addressing (stable UUID + human
// key), forward reference edges with derived backlinks, tags, revisions and an
// activity timeline. It mirrors the sibling tasker engine's shape; storage and
// HTTP transport wrap this package (later increments: internal/store,
// internal/httpapi). Design follows docs/knowledge-base-canon.md.
package kb

import (
	"errors"
	"regexp"
)

// Error classes returned by the domain; the HTTP layer maps them to status
// codes (404/400/409). Wrap with %w to add detail.
var (
	// ErrNotFound — a requested entity does not exist.
	ErrNotFound = errors.New("not found")
	// ErrValidation — input failed a domain invariant.
	ErrValidation = errors.New("validation")
	// ErrConflict — the operation conflicts with current state (e.g. duplicate key).
	ErrConflict = errors.New("conflict")
)

// keyPattern matches a stable node/workspace key such as "KNOW-12": an
// UPPER-case workspace prefix, a hyphen, then a monotonic sequence number.
var keyPattern = regexp.MustCompile(`^[A-Z][A-Z0-9]*-\d+$`)

// LooksLikeKey reports whether idOrKey is a stable human key ("WS-n") rather
// than a UUID. Dual-id resolution tries the UUID first, then falls back to the
// key; the two formats do not overlap, so the probe order is unambiguous.
func LooksLikeKey(idOrKey string) bool {
	return keyPattern.MatchString(idOrKey)
}
