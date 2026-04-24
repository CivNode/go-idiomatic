package goidiomatic

import (
	"github.com/CivNode/go-idiomatic/rules"
)

// Severity grades how strong a finding is.
type Severity = rules.Severity

// Severity levels re-exported from the rules package so callers can write
// goidiomatic.Error etc. without importing the subpackage directly.
const (
	Info  = rules.Info
	Warn  = rules.Warn
	Error = rules.Error
)

// Rule is a single pedagogical check.
type Rule = rules.Rule

// Finding is one hit on a Rule.
type Finding = rules.Finding

// Fix is an optional suggested edit attached to a Finding.
type Fix = rules.Fix

// TextEdit is a single replacement in source.
type TextEdit = rules.TextEdit
