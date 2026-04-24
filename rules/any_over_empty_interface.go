package rules

import (
	"go/ast"
	"go/build"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"

	goidiomatic "github.com/CivNode/go-idiomatic"
)

// AnyOverEmptyInterface flags uses of interface{} and suggests any, but
// only when the build context's language version is Go 1.18 or newer. On
// older targets the suggestion would be a compile error, so the rule stays
// silent.
var AnyOverEmptyInterface goidiomatic.Rule = anyOverEmptyInterface{}

type anyOverEmptyInterface struct{}

func (anyOverEmptyInterface) ID() string   { return "any-over-empty-interface" }
func (anyOverEmptyInterface) Name() string { return "prefer any over interface{}" }
func (anyOverEmptyInterface) Description() string {
	return "Since Go 1.18 the predeclared alias any reads more cleanly than interface{}. This rule only fires when the module language version is 1.18 or newer."
}
func (anyOverEmptyInterface) Severity() goidiomatic.Severity { return goidiomatic.Info }

func (r anyOverEmptyInterface) Check(pass *analysis.Pass) ([]goidiomatic.Finding, error) {
	if !goLangAtLeast(1, 18) {
		return nil, nil
	}

	var out []goidiomatic.Finding
	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			it, ok := n.(*ast.InterfaceType)
			if !ok {
				return true
			}
			if it.Methods == nil || len(it.Methods.List) != 0 {
				return true
			}
			out = append(out, goidiomatic.Finding{
				RuleID:   r.ID(),
				Message:  "use any instead of interface{} on Go 1.18+",
				Pos:      pass.Fset.Position(it.Pos()),
				Severity: r.Severity(),
			})
			return true
		})
	}
	return out, nil
}

// goLangAtLeast reports whether the current build language is at least the
// given major.minor. It consults build.Default.ReleaseTags, which is the
// canonical way to learn the target Go version in a package-loading context.
func goLangAtLeast(major, minor int) bool {
	tags := build.Default.ReleaseTags
	want := "go" + strconv.Itoa(major) + "." + strconv.Itoa(minor)
	for _, t := range tags {
		if t == want {
			return true
		}
	}
	// If the exact tag is missing but a newer tag exists, still true.
	for _, t := range tags {
		if !strings.HasPrefix(t, "go") {
			continue
		}
		v := strings.TrimPrefix(t, "go")
		parts := strings.SplitN(v, ".", 2)
		if len(parts) != 2 {
			continue
		}
		maj, err1 := strconv.Atoi(parts[0])
		min, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			continue
		}
		if maj > major || (maj == major && min >= minor) {
			return true
		}
	}
	return false
}
