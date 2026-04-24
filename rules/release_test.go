package rules

import "testing"

// These helpers live in an internal test (package rules, not rules_test)
// so we can reach releaseAtLeast and isErrorTyped / isContextType without
// reshuffling exports.

func TestReleaseAtLeast(t *testing.T) {
	cases := []struct {
		name string
		tags []string
		maj  int
		min  int
		want bool
	}{
		{"exact match", []string{"go1.17", "go1.18"}, 1, 18, true},
		{"newer tag covers it", []string{"go1.20"}, 1, 18, true},
		{"only older tags", []string{"go1.16", "go1.17"}, 1, 18, false},
		{"major bump", []string{"go2.0"}, 1, 30, true},
		{"ignores junk tags", []string{"linux", "go1.18"}, 1, 18, true},
		{"mis-shaped tag tolerated", []string{"gofoo", "go1.18"}, 1, 18, true},
		{"malformed version ignored", []string{"go1", "go1.a", "go1.18"}, 1, 18, true},
		{"empty tags", nil, 1, 18, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := releaseAtLeast(c.tags, c.maj, c.min); got != c.want {
				t.Errorf("releaseAtLeast(%v, %d, %d) = %v want %v", c.tags, c.maj, c.min, got, c.want)
			}
		})
	}
}

func TestLastSegment(t *testing.T) {
	cases := map[string]string{
		"context":           "context",
		"github.com/a/b":    "b",
		"github.com/a/b/c":  "c",
		"":                  "",
		"/trailing":         "trailing",
		"no-slash-but-dash": "no-slash-but-dash",
	}
	for in, want := range cases {
		if got := lastSegment(in); got != want {
			t.Errorf("lastSegment(%q) = %q want %q", in, got, want)
		}
	}
}

func TestAnalyzerName(t *testing.T) {
	cases := map[string]string{
		"simple":             "simple",
		"prefer-range-int":   "prefer_range_int",
		"already_underscore": "already_underscore",
		"":                   "",
	}
	for in, want := range cases {
		if got := analyzerName(in); got != want {
			t.Errorf("analyzerName(%q) = %q want %q", in, got, want)
		}
	}
}
