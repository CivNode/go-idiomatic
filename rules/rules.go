// Package rules implements the built-in pedagogical rules for go-idiomatic.
//
// Every rule is a value exposing the goidiomatic.Rule interface and also
// exposes an *analysis.Analyzer via Analyzer / Analyzers so it can be
// plugged into singlechecker, multichecker, or analysistest.Run.
package rules

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"

	goidiomatic "github.com/CivNode/go-idiomatic"
)

// All returns every built-in rule in a stable order.
func All() []goidiomatic.Rule {
	return []goidiomatic.Rule{
		PreferRangeInt,
		ErrorsIsAs,
		AnyOverEmptyInterface,
		ContextFirstArg,
		NoSleepForCoordination,
	}
}

// Analyzer wraps every built-in rule as a single Analyzer suitable for
// singlechecker / multichecker.
var Analyzer = newAnalyzer("goidiomatic", "pedagogical Go-idiom checks beyond golangci-lint", All())

// Analyzers returns one analyzer per rule, keyed by rule ID. Useful when a
// caller wants to surface findings per-rule or run analysistest on a single
// rule.
func Analyzers() map[string]*analysis.Analyzer {
	out := make(map[string]*analysis.Analyzer, len(All()))
	for _, r := range All() {
		out[r.ID()] = newAnalyzer(r.ID(), r.Description(), []goidiomatic.Rule{r})
	}
	return out
}

func newAnalyzer(name, doc string, rs []goidiomatic.Rule) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: analyzerName(name),
		Doc:  doc,
		Run: func(pass *analysis.Pass) (any, error) {
			for _, r := range rs {
				findings, err := r.Check(pass)
				if err != nil {
					return nil, err
				}
				for _, f := range findings {
					pos := positionToPos(pass, f.Pos)
					pass.Report(analysis.Diagnostic{
						Pos:      pos,
						Category: r.ID(),
						Message:  f.Message,
					})
				}
			}
			return nil, nil
		},
	}
}

// analyzerName converts a dashed rule id into a valid analyzer name. The
// analysis framework requires the name to be a valid Go identifier.
func analyzerName(id string) string {
	out := make([]byte, 0, len(id))
	for i := 0; i < len(id); i++ {
		c := id[i]
		if c == '-' {
			out = append(out, '_')
		} else {
			out = append(out, c)
		}
	}
	return string(out)
}

// positionToPos converts a token.Position back to a token.Pos within pass.
func positionToPos(pass *analysis.Pass, p token.Position) token.Pos {
	var result token.Pos
	pass.Fset.Iterate(func(f *token.File) bool {
		if f.Name() != p.Filename {
			return true
		}
		if p.Line < 1 || p.Line > f.LineCount() {
			return false
		}
		result = f.LineStart(p.Line) + token.Pos(p.Column-1)
		return false
	})
	return result
}

// Run parses src as a single Go file in a synthetic package and runs the
// given rule against it. It is a test-oriented helper that does not require
// module resolution or GOPATH.
func Run(rule goidiomatic.Rule, src []byte) ([]goidiomatic.Finding, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "src.go", src, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	conf := types.Config{Importer: nopImporter{}, Error: func(error) {}}
	info := &types.Info{
		Types:      map[ast.Expr]types.TypeAndValue{},
		Defs:       map[*ast.Ident]types.Object{},
		Uses:       map[*ast.Ident]types.Object{},
		Selections: map[*ast.SelectorExpr]*types.Selection{},
	}
	pkg, _ := conf.Check("x", fset, []*ast.File{file}, info)

	pass := &analysis.Pass{
		Analyzer:  Analyzer,
		Fset:      fset,
		Files:     []*ast.File{file},
		Pkg:       pkg,
		TypesInfo: info,
		Report:    func(analysis.Diagnostic) {},
	}
	return rule.Check(pass)
}

// nopImporter is a types.Importer that returns a synthetic package for any
// import path. Rules that need real type information should rely on
// syntactic shape rather than full resolution so they work in both modes.
type nopImporter struct{}

func (nopImporter) Import(path string) (*types.Package, error) {
	return types.NewPackage(path, lastSegment(path)), nil
}

func lastSegment(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			return path[i+1:]
		}
	}
	return path
}
