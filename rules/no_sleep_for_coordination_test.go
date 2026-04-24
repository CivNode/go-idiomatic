package rules_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/CivNode/go-idiomatic/rules"
)

func TestNoSleepForCoordination_AnalysisTest(t *testing.T) {
	a := rules.Analyzers()[rules.NoSleepForCoordination.ID()]
	dir := analysistest.TestData()
	analysistest.Run(t, dir, a, "nosleep_ok", "nosleep_bad")
}

func TestNoSleepForCoordination_FlagsChannel(t *testing.T) {
	src := []byte(`package x
import "time"
func F(done <-chan struct{}) { time.Sleep(time.Second); <-done }`)
	got, err := rules.Run(rules.NoSleepForCoordination, src)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("want 1, got %d: %+v", len(got), got)
	}
}

func TestNoSleepForCoordination_FlagsContext(t *testing.T) {
	src := []byte(`package x
import (
	"context"
	"time"
)
func F(ctx context.Context) { time.Sleep(time.Second); _ = ctx }`)
	got, err := rules.Run(rules.NoSleepForCoordination, src)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("want 1, got %d: %+v", len(got), got)
	}
}

func TestNoSleepForCoordination_IgnoresStandaloneSleep(t *testing.T) {
	src := []byte(`package x
import "time"
func F() { time.Sleep(time.Second) }`)
	got, _ := rules.Run(rules.NoSleepForCoordination, src)
	if len(got) != 0 {
		t.Fatalf("want 0, got %d", len(got))
	}
}
