package goidiomatic

import (
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// Severity grades how strong a finding is.
type Severity int

// Severity levels, ordered from weakest to strongest.
const (
	// Info is a stylistic nudge.
	Info Severity = iota + 1
	// Warn is a likely issue worth addressing.
	Warn
	// Error is almost certainly a bug or anti-pattern.
	Error
)

// String returns the human label for a severity.
func (s Severity) String() string {
	switch s {
	case Info:
		return "info"
	case Warn:
		return "warn"
	case Error:
		return "error"
	default:
		return "unknown"
	}
}

// Rule is a single pedagogical check. Each rule runs over an
// analysis.Pass and returns findings without mutating it.
type Rule interface {
	ID() string
	Name() string
	Description() string
	Severity() Severity
	Check(pass *analysis.Pass) ([]Finding, error)
}

// Finding is one hit on a Rule.
type Finding struct {
	RuleID   string
	Message  string
	Pos      token.Position
	Severity Severity
	Fix      *Fix
}

// Fix is an optional suggested edit attached to a Finding.
type Fix struct {
	Description string
	TextEdits   []TextEdit
}

// TextEdit is a single replacement in source.
type TextEdit struct {
	Start, End token.Position
	NewText    string
}
