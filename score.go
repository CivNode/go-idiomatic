package goidiomatic

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// Score parses src, runs each rule against it, and returns an aggregate
// score on [0, 100] plus all findings. The score starts at 100 and deducts
// 15 per Error, 7 per Warn, 2 per Info, floored at 0.
//
// Score is deliberately single-file and loader-free so callers can grade
// snippets straight from memory. For package-level analysis use the
// analyzer exposed by the rules package with the standard analysis driver
// instead.
func Score(src []byte, ruleset []Rule) (int, []Finding, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "src.go", src, parser.AllErrors)
	if err != nil {
		return 0, nil, fmt.Errorf("parse: %w", err)
	}

	conf := types.Config{Importer: syntheticImporter{}, Error: func(error) {}}
	info := &types.Info{
		Types:      map[ast.Expr]types.TypeAndValue{},
		Defs:       map[*ast.Ident]types.Object{},
		Uses:       map[*ast.Ident]types.Object{},
		Selections: map[*ast.SelectorExpr]*types.Selection{},
	}
	pkg, _ := conf.Check("x", fset, []*ast.File{file}, info)

	pass := &analysis.Pass{
		Analyzer:  &analysis.Analyzer{Name: "score", Doc: "score pass"},
		Fset:      fset,
		Files:     []*ast.File{file},
		Pkg:       pkg,
		TypesInfo: info,
		Report:    func(analysis.Diagnostic) {},
	}

	var all []Finding
	for _, r := range ruleset {
		findings, err := r.Check(pass)
		if err != nil {
			return 0, nil, fmt.Errorf("%s: %w", r.ID(), err)
		}
		all = append(all, findings...)
	}

	score := 100
	for _, f := range all {
		switch f.Severity {
		case Error:
			score -= 15
		case Warn:
			score -= 7
		case Info:
			score -= 2
		}
	}
	if score < 0 {
		score = 0
	}
	return score, all, nil
}

// syntheticImporter satisfies go/types.Importer with stub packages for any
// import path so Score can type-check snippets without a module on disk.
type syntheticImporter struct{}

func (syntheticImporter) Import(path string) (*types.Package, error) {
	return types.NewPackage(path, shortName(path)), nil
}

func shortName(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			return path[i+1:]
		}
	}
	return path
}
