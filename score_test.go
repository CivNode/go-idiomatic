package goidiomatic_test

import (
	"testing"

	goidiomatic "github.com/CivNode/go-idiomatic"
)

func TestScore_PerfectSourceScoresHundred(t *testing.T) {
	src := []byte(`package x
func F(xs []int) int {
	sum := 0
	for _, v := range xs {
		sum += v
	}
	return sum
}`)
	score, findings, err := goidiomatic.Score(src, goidiomatic.DefaultRules())
	if err != nil {
		t.Fatalf("score: %v", err)
	}
	if score != 100 {
		t.Fatalf("want 100, got %d: %+v", score, findings)
	}
	if len(findings) != 0 {
		t.Fatalf("want 0 findings, got %d", len(findings))
	}
}

func TestScore_DeductsPerSeverity(t *testing.T) {
	// Info severity: deducts 2.
	src := []byte(`package x
func F(xs []int) int {
	sum := 0
	for i := 0; i < len(xs); i++ { sum += xs[i] }
	return sum
}`)
	score, findings, err := goidiomatic.Score(src, goidiomatic.DefaultRules())
	if err != nil {
		t.Fatalf("score: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("want 1 finding, got %d", len(findings))
	}
	if score != 98 {
		t.Fatalf("want 98, got %d", score)
	}
}

func TestScore_FloorAtZero(t *testing.T) {
	// Ten no-sleep-for-coordination Error findings (15 * 10 = 150) > 100.
	src := []byte(`package x
import (
	"context"
	"time"
)
func A(ctx context.Context) { time.Sleep(time.Second); _ = ctx }
func B(ctx context.Context) { time.Sleep(time.Second); _ = ctx }
func C(ctx context.Context) { time.Sleep(time.Second); _ = ctx }
func D(ctx context.Context) { time.Sleep(time.Second); _ = ctx }
func E(ctx context.Context) { time.Sleep(time.Second); _ = ctx }
func F(ctx context.Context) { time.Sleep(time.Second); _ = ctx }
func G(ctx context.Context) { time.Sleep(time.Second); _ = ctx }
`)
	score, _, err := goidiomatic.Score(src, goidiomatic.DefaultRules())
	if err != nil {
		t.Fatal(err)
	}
	if score != 0 {
		t.Fatalf("want 0, got %d", score)
	}
}

func TestScore_InvalidSourceReturnsError(t *testing.T) {
	_, _, err := goidiomatic.Score([]byte("not go"), goidiomatic.DefaultRules())
	if err == nil {
		t.Fatal("want parse error, got nil")
	}
}

func TestDefaultRules_StableOrdering(t *testing.T) {
	a := goidiomatic.DefaultRules()
	b := goidiomatic.DefaultRules()
	if len(a) != len(b) {
		t.Fatalf("length drift")
	}
	for i := range a {
		if a[i].ID() != b[i].ID() {
			t.Fatalf("order drift at %d: %s vs %s", i, a[i].ID(), b[i].ID())
		}
	}
	if len(a) < 5 {
		t.Fatalf("want at least 5 rules, got %d", len(a))
	}
}

func TestSeverity_String(t *testing.T) {
	cases := map[goidiomatic.Severity]string{
		goidiomatic.Info:            "info",
		goidiomatic.Warn:            "warn",
		goidiomatic.Error:           "error",
		goidiomatic.Severity(99999): "unknown",
	}
	for s, want := range cases {
		if got := s.String(); got != want {
			t.Errorf("%d: got %q want %q", s, got, want)
		}
	}
}
