package goidiomatic

import "github.com/CivNode/go-idiomatic/rules"

// DefaultRules returns the built-in ruleset in a stable order. The returned
// slice is a fresh copy; callers may reorder or extend it without
// affecting future calls.
func DefaultRules() []Rule {
	all := rules.All()
	out := make([]Rule, len(all))
	copy(out, all)
	return out
}
