// Package goidiomatic is a pedagogical Go-idiom linter.
//
// Rules live under the rules/ subpackage. Each rule implements the Rule
// interface and ships with golden fixtures under rules/testdata/<id>/
// (ok.go for code that must not flag, bad.go for code that must flag).
//
// The Score function aggregates rule findings into a 0..100 score,
// deducting per severity level.
//
// See https://github.com/CivNode/go-idiomatic for details.
package goidiomatic
