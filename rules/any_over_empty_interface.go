package rules

import (
	"go/ast"
	"go/build"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// AnyOverEmptyInterface flags uses of interface{} and suggests any, but
// only when the build context's language version is Go 1.18 or newer. On
// older targets the suggestion would be a compile error, so the rule stays
// silent.
var AnyOverEmptyInterface Rule = anyOverEmptyInterface{}

type anyOverEmptyInterface struct{}

func (anyOverEmptyInterface) ID() string   { return "any-over-empty-interface" }
func (anyOverEmptyInterface) Name() string { return "prefer any over interface{}" }
func (anyOverEmptyInterface) Description() string {
	return "Since Go 1.18 the predeclared alias any reads more cleanly than interface{}. This rule only fires when the module language version is 1.18 or newer."
}
func (anyOverEmptyInterface) Severity() Severity { return Info }

func (r anyOverEmptyInterface) Check(pass *analysis.Pass) ([]Finding, error) {
	if !releaseAtLeast(build.Default.ReleaseTags, 1, 18) {
		return nil, nil
	}

	var out []Finding
	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			it, ok := n.(*ast.InterfaceType)
			if !ok {
				return true
			}
			if it.Methods == nil || len(it.Methods.List) != 0 {
				return true
			}
			out = append(out, Finding{
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

// releaseAtLeast reports whether the given release tags contain a Go
// version equal to or newer than major.minor. Extracted so tests can pass
// synthetic tag lists in, rather than depending on the host toolchain.
func releaseAtLeast(tags []string, major, minor int) bool {
	want := "go" + strconv.Itoa(major) + "." + strconv.Itoa(minor)
	for _, t := range tags {
		if t == want {
			return true
		}
	}
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
