package rules_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/CivNode/go-idiomatic/rules"
)

func TestPreferRangeInt_AnalysisTest(t *testing.T) {
	a := rules.Analyzers()[rules.PreferRangeInt.ID()]
	dir := analysistest.TestData()
	analysistest.Run(t, dir, a, "preferrangeint_ok", "preferrangeint_bad")
}

func TestPreferRangeInt_FlagsClassic(t *testing.T) {
	src := []byte(`package x
func F(xs []int) int {
	sum := 0
	for i := 0; i < len(xs); i++ { sum += xs[i] }
	return sum
}`)
	got, err := rules.Run(rules.PreferRangeInt, src)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("want 1 finding, got %d: %+v", len(got), got)
	}
	if got[0].RuleID != "prefer-range-int" {
		t.Fatalf("rule id %q", got[0].RuleID)
	}
}

func TestPreferRangeInt_IgnoresCountdown(t *testing.T) {
	src := []byte(`package x
func F(xs []int) { for i := len(xs)-1; i >= 0; i-- { _ = xs[i] } }`)
	got, err := rules.Run(rules.PreferRangeInt, src)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("want 0 findings, got %d", len(got))
	}
}

func TestPreferRangeInt_IgnoresBiggerStep(t *testing.T) {
	src := []byte(`package x
func F(xs []int) { for i := 0; i < len(xs); i += 2 { _ = xs[i] } }`)
	got, _ := rules.Run(rules.PreferRangeInt, src)
	if len(got) != 0 {
		t.Fatalf("want 0 findings, got %d", len(got))
	}
}
