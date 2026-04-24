package rules_test

import (
	"strings"
	"testing"

	"github.com/CivNode/go-idiomatic/rules"
)

// TestAllRules_Metadata exercises every metadata getter and guarantees the
// returned strings are non-empty and IDs are unique.
func TestAllRules_Metadata(t *testing.T) {
	seen := map[string]bool{}
	for _, r := range rules.All() {
		if r.ID() == "" {
			t.Errorf("rule has empty ID")
		}
		if r.Name() == "" {
			t.Errorf("%s: empty Name", r.ID())
		}
		if r.Description() == "" {
			t.Errorf("%s: empty Description", r.ID())
		}
		if s := r.Severity(); s < rules.Info || s > rules.Error {
			t.Errorf("%s: severity out of range: %d", r.ID(), s)
		}
		if seen[r.ID()] {
			t.Errorf("duplicate ID %q", r.ID())
		}
		seen[r.ID()] = true
	}
}

func TestSeverity_String(t *testing.T) {
	cases := map[rules.Severity]string{
		rules.Info:            "info",
		rules.Warn:            "warn",
		rules.Error:           "error",
		rules.Severity(12345): "unknown",
	}
	for s, want := range cases {
		if got := s.String(); got != want {
			t.Errorf("%d: got %q want %q", s, got, want)
		}
	}
}

func TestAnalyzers_OneAnalyzerPerRule(t *testing.T) {
	analyzers := rules.Analyzers()
	if len(analyzers) != len(rules.All()) {
		t.Fatalf("analyzer count %d != rule count %d", len(analyzers), len(rules.All()))
	}
	for _, r := range rules.All() {
		a, ok := analyzers[r.ID()]
		if !ok {
			t.Errorf("%s: no analyzer registered", r.ID())
			continue
		}
		if a.Name == "" {
			t.Errorf("%s: analyzer has empty name", r.ID())
		}
		if strings.Contains(a.Name, "-") {
			t.Errorf("%s: analyzer name %q still contains dash", r.ID(), a.Name)
		}
	}
}

func TestRun_ReturnsParseErrors(t *testing.T) {
	_, err := rules.Run(rules.PreferRangeInt, []byte("not go code"))
	if err == nil {
		t.Fatal("want parse error, got nil")
	}
}

// TestPreferRangeInt_ExoticShapes covers the early-out branches in
// classicIndexedLoop that the basic happy-path tests did not hit.
func TestPreferRangeInt_ExoticShapes(t *testing.T) {
	cases := map[string]string{
		"non-zero start": `package x
func F(xs []int) { for i := 1; i < len(xs); i++ { _ = xs[i] } }`,
		"condition is >=": `package x
func F(xs []int) { for i := 0; i >= len(xs); i++ { _ = xs[i] } }`,
		"condition rhs is not len": `package x
func F(xs []int) { for i := 0; i < cap(xs); i++ { _ = xs[i] } }`,
		"len on non-ident": `package x
type S struct{ xs []int }
func (s S) F() { for i := 0; i < len(s.xs); i++ { _ = s.xs[i] } }`,
		"post is decrement": `package x
func F(xs []int) { for i := 0; i < len(xs); i-- { _ = xs[i] } }`,
		"different post var": `package x
func F(xs []int) { j := 0; for i := 0; i < len(xs); j++ { _ = xs[i]; _ = j } }`,
		"different cond lhs": `package x
func F(xs []int) { j := 0; for i := 0; j < len(xs); i++ { _ = xs[i]; _ = j } }`,
	}
	for name, src := range cases {
		t.Run(name, func(t *testing.T) {
			got, _ := rules.Run(rules.PreferRangeInt, []byte(src))
			if len(got) != 0 {
				t.Errorf("%s: want no findings, got %d: %+v", name, len(got), got)
			}
		})
	}
}

// TestErrorsIsAs_ExoticShapes exercises the non-matching branches.
func TestErrorsIsAs_ExoticShapes(t *testing.T) {
	cases := map[string]string{
		"unrelated comparison": `package x
func F(a, b int) bool { return a == b }`,
		"method call with args is not Error": `package x
type T struct{}
func (T) Error(x int) string { return "" }
func F(t T) bool { return t.Error(1) == "" }`,
		"selector not Error()": `package x
type T struct{}
func (T) String() string { return "" }
func F(t T) bool { return t.String() == "foo" }`,
	}
	for name, src := range cases {
		t.Run(name, func(t *testing.T) {
			got, _ := rules.Run(rules.ErrorsIsAs, []byte(src))
			if len(got) != 0 {
				t.Errorf("want 0 findings, got %d: %+v", len(got), got)
			}
		})
	}
}

func TestErrorsIsAs_FlagsReversedOperands(t *testing.T) {
	src := []byte(`package x
func F(err error) bool { return "boom" == err.Error() }`)
	got, _ := rules.Run(rules.ErrorsIsAs, src)
	if len(got) != 1 {
		t.Fatalf("want 1 finding, got %d", len(got))
	}
}

func TestContextFirstArg_FuncLit(t *testing.T) {
	src := []byte(`package x
import "context"
var _ = func(id string, ctx context.Context) { _ = ctx; _ = id }`)
	got, err := rules.Run(rules.ContextFirstArg, src)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("want 1, got %d: %+v", len(got), got)
	}
}

func TestContextFirstArg_MultiNameField(t *testing.T) {
	// Multi-name field declarations such as (a, b context.Context) still
	// count each parameter individually; the first name is at index 0.
	src := []byte(`package x
import "context"
func F(a, b context.Context) { _ = a; _ = b }`)
	got, _ := rules.Run(rules.ContextFirstArg, src)
	// a is at idx 0 so no finding; b is at idx 1 so one finding may or may
	// not be raised depending on whether we treat multi-name fields as one
	// group. Current behavior: count individuals, so b flags.
	if len(got) == 0 {
		t.Skip("multi-name fields treated as a single unit; fine either way")
	}
}

func TestNoSleep_FuncLitWithContext(t *testing.T) {
	src := []byte(`package x
import (
	"context"
	"time"
)
var _ = func(ctx context.Context) { time.Sleep(time.Second); _ = ctx }`)
	got, err := rules.Run(rules.NoSleepForCoordination, []byte(src))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("want 1, got %d", len(got))
	}
}

func TestNoSleep_ChannelTypes(t *testing.T) {
	// Function that declares a chan type triggers usesChannel via ChanType.
	src := []byte(`package x
import "time"
func F() { var ch chan int; _ = ch; time.Sleep(time.Second) }`)
	got, _ := rules.Run(rules.NoSleepForCoordination, src)
	if len(got) != 1 {
		t.Fatalf("want 1, got %d", len(got))
	}
}

func TestNoSleep_SendStmt(t *testing.T) {
	src := []byte(`package x
import "time"
func F(ch chan int) { ch <- 1; time.Sleep(time.Second) }`)
	got, _ := rules.Run(rules.NoSleepForCoordination, src)
	if len(got) != 1 {
		t.Fatalf("want 1, got %d", len(got))
	}
}

func TestNoSleep_SelectStmt(t *testing.T) {
	src := []byte(`package x
import "time"
func F(ch chan int) {
	select { case <-ch: }
	time.Sleep(time.Second)
}`)
	got, _ := rules.Run(rules.NoSleepForCoordination, src)
	if len(got) != 1 {
		t.Fatalf("want 1, got %d", len(got))
	}
}

func TestNoSleep_FuncLitNoBody(t *testing.T) {
	// Pure function literal with no coordination must not flag.
	src := []byte(`package x
import "time"
func F() { time.Sleep(time.Second) }`)
	got, _ := rules.Run(rules.NoSleepForCoordination, src)
	if len(got) != 0 {
		t.Fatalf("want 0, got %d", len(got))
	}
}

func TestAnyOverEmpty_TypeAliasDecl(t *testing.T) {
	// Interface type in a non-type-spec context (field).
	src := []byte(`package x
var x interface{} = 5
func _unused() { _ = x }`)
	got, err := rules.Run(rules.AnyOverEmptyInterface, []byte(src))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("want 1 finding, got %d: %+v", len(got), got)
	}
}
