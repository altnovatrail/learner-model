package learnermodel

import (
	"errors"
	"strconv"
)

// The kinds of error the library returns, as values and never by panicking
// (docs/architecture.md, section 3.7). errors.Is tells them apart; the message says where.
var (
	// ErrInvalid is an input that breaks a rule the library relies on: a duplicate ID, a place
	// outside the bands, a skill the graph does not know.
	ErrInvalid = errors.New("learnermodel: invalid input")
	// ErrCorrupt is bytes that are no encoding of evidence or a state.
	ErrCorrupt = errors.New("learnermodel: corrupt bytes")
	// ErrNotFinite is a NaN or an infinity, which the encoding refuses both ways.
	ErrNotFinite = errors.New("learnermodel: a number that is not finite")
)

// detailError is one of the kinds above with what and where.
type detailError struct {
	kind   error
	detail string
}

func (e *detailError) Error() string { return e.kind.Error() + ": " + e.detail }

func (e *detailError) Unwrap() error { return e.kind }

func invalid(detail string) error { return &detailError{ErrInvalid, detail} }

func itoa(v int32) string { return strconv.Itoa(int(v)) }
