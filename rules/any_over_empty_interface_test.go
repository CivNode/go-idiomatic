package rules_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/CivNode/go-idiomatic/rules"
)

func TestAnyOverEmptyInterface_AnalysisTest(t *testing.T) {
	a := rules.Analyzers()[rules.AnyOverEmptyInterface.ID()]
	dir := analysistest.TestData()
	analysistest.Run(t, dir, a, "anyoverempty_ok", "anyoverempty_bad")
}

func TestAnyOverEmptyInterface_Flags(t *testing.T) {
	src := []byte(`package x
type Box struct { V interface{} }
func F(v interface{}) interface{} { return v }`)
	got, err := rules.Run(rules.AnyOverEmptyInterface, src)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 {
		t.Fatalf("want 3 findings, got %d: %+v", len(got), got)
	}
}

func TestAnyOverEmptyInterface_IgnoresNonEmpty(t *testing.T) {
	src := []byte(`package x
type S interface { String() string }`)
	got, _ := rules.Run(rules.AnyOverEmptyInterface, src)
	if len(got) != 0 {
		t.Fatalf("want 0, got %d: %+v", len(got), got)
	}
}
