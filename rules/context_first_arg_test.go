package rules_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/CivNode/go-idiomatic/rules"
)

func TestContextFirstArg_AnalysisTest(t *testing.T) {
	a := rules.Analyzers()[rules.ContextFirstArg.ID()]
	dir := analysistest.TestData()
	analysistest.Run(t, dir, a, "ctxfirst_ok", "ctxfirst_bad")
}

func TestContextFirstArg_Flags(t *testing.T) {
	src := []byte(`package x
import "context"
func F(id string, ctx context.Context) { _ = ctx; _ = id }`)
	got, err := rules.Run(rules.ContextFirstArg, src)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("want 1, got %d: %+v", len(got), got)
	}
}

func TestContextFirstArg_IgnoresContextFirst(t *testing.T) {
	src := []byte(`package x
import "context"
func F(ctx context.Context, id string) { _ = ctx; _ = id }`)
	got, _ := rules.Run(rules.ContextFirstArg, src)
	if len(got) != 0 {
		t.Fatalf("want 0, got %d", len(got))
	}
}

func TestContextFirstArg_SingleParamFuncs(t *testing.T) {
	src := []byte(`package x
import "context"
func F(ctx context.Context) { _ = ctx }`)
	got, _ := rules.Run(rules.ContextFirstArg, src)
	if len(got) != 0 {
		t.Fatalf("want 0, got %d", len(got))
	}
}
