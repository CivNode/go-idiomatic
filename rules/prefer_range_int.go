package rules

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"

	goidiomatic "github.com/CivNode/go-idiomatic"
)

// PreferRangeInt flags classic indexed loops that could be rewritten with
// range, for example for i := 0; i < len(xs); i++. It ignores countdowns
// and steps other than i++.
var PreferRangeInt goidiomatic.Rule = preferRangeInt{}

type preferRangeInt struct{}

func (preferRangeInt) ID() string   { return "prefer-range-int" }
func (preferRangeInt) Name() string { return "prefer range over indexed for" }
func (preferRangeInt) Description() string {
	return "Classic for i := 0; i < len(x); i++ loops can be written as for i := range x, which is shorter and harder to misuse."
}
func (preferRangeInt) Severity() goidiomatic.Severity { return goidiomatic.Info }

func (r preferRangeInt) Check(pass *analysis.Pass) ([]goidiomatic.Finding, error) {
	var out []goidiomatic.Finding
	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			loop, ok := n.(*ast.ForStmt)
			if !ok || loop.Init == nil || loop.Cond == nil || loop.Post == nil {
				return true
			}
			target, ok := classicIndexedLoop(loop)
			if !ok {
				return true
			}
			pos := pass.Fset.Position(loop.Pos())
			out = append(out, goidiomatic.Finding{
				RuleID:   r.ID(),
				Message:  fmt.Sprintf("prefer `for i := range %s` over classic indexed loop", target),
				Pos:      pos,
				Severity: r.Severity(),
			})
			return true
		})
	}
	return out, nil
}

// classicIndexedLoop returns the target name and true when loop is the shape
// for i := 0; i < len(target); i++ (or <=, with post increment). It refuses
// countdowns and compound post statements.
func classicIndexedLoop(loop *ast.ForStmt) (string, bool) {
	// Init must be `i := 0`.
	assign, ok := loop.Init.(*ast.AssignStmt)
	if !ok || assign.Tok != token.DEFINE || len(assign.Lhs) != 1 || len(assign.Rhs) != 1 {
		return "", false
	}
	idxIdent, ok := assign.Lhs[0].(*ast.Ident)
	if !ok {
		return "", false
	}
	zero, ok := assign.Rhs[0].(*ast.BasicLit)
	if !ok || zero.Kind != token.INT || zero.Value != "0" {
		return "", false
	}

	// Cond must be `i < len(target)`.
	bin, ok := loop.Cond.(*ast.BinaryExpr)
	if !ok || bin.Op != token.LSS {
		return "", false
	}
	left, ok := bin.X.(*ast.Ident)
	if !ok || left.Name != idxIdent.Name {
		return "", false
	}
	call, ok := bin.Y.(*ast.CallExpr)
	if !ok || len(call.Args) != 1 {
		return "", false
	}
	fnIdent, ok := call.Fun.(*ast.Ident)
	if !ok || fnIdent.Name != "len" {
		return "", false
	}
	argIdent, ok := call.Args[0].(*ast.Ident)
	if !ok {
		return "", false
	}

	// Post must be `i++`.
	inc, ok := loop.Post.(*ast.IncDecStmt)
	if !ok || inc.Tok != token.INC {
		return "", false
	}
	postIdent, ok := inc.X.(*ast.Ident)
	if !ok || postIdent.Name != idxIdent.Name {
		return "", false
	}

	return argIdent.Name, true
}
