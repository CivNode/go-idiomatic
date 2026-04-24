package rules_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/CivNode/go-idiomatic/rules"
)

func TestErrorsIsAs_AnalysisTest(t *testing.T) {
	a := rules.Analyzers()[rules.ErrorsIsAs.ID()]
	dir := analysistest.TestData()
	analysistest.Run(t, dir, a, "errorsisas_ok", "errorsisas_bad")
}

func TestErrorsIsAs_FlagsStringCompare(t *testing.T) {
	src := []byte(`package x
func F(err error) bool { return err.Error() == "boom" }`)
	got, err := rules.Run(rules.ErrorsIsAs, src)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("want 1, got %d: %+v", len(got), got)
	}
}

func TestErrorsIsAs_IgnoresNilCheck(t *testing.T) {
	src := []byte(`package x
func F(err error) bool { return err == nil }`)
	got, _ := rules.Run(rules.ErrorsIsAs, src)
	if len(got) != 0 {
		t.Fatalf("want 0, got %d", len(got))
	}
}
