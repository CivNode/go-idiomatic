package rules

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"

	goidiomatic "github.com/CivNode/go-idiomatic"
)

// ErrorsIsAs flags two common anti-patterns:
//
//  1. comparing err.Error() to a literal string
//  2. comparing an error value to a sentinel with ==
//
// Both should use errors.Is (or errors.As for type switches) so wrapped
// errors still match. The check is structural, augmented with type info
// when it is available.
var ErrorsIsAs goidiomatic.Rule = errorsIsAs{}

type errorsIsAs struct{}

func (errorsIsAs) ID() string   { return "errors-is-as" }
func (errorsIsAs) Name() string { return "prefer errors.Is / errors.As" }
func (errorsIsAs) Description() string {
	return "Comparing err.Error() to a string, or an error value to a sentinel with ==, misses wrapped errors. Use errors.Is or errors.As."
}
func (errorsIsAs) Severity() goidiomatic.Severity { return goidiomatic.Warn }

func (r errorsIsAs) Check(pass *analysis.Pass) ([]goidiomatic.Finding, error) {
	var out []goidiomatic.Finding
	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			bin, ok := n.(*ast.BinaryExpr)
			if !ok || (bin.Op != token.EQL && bin.Op != token.NEQ) {
				return true
			}

			if msg, hit := errorStringCompare(bin); hit {
				out = append(out, goidiomatic.Finding{
					RuleID:   r.ID(),
					Message:  msg,
					Pos:      pass.Fset.Position(bin.Pos()),
					Severity: r.Severity(),
				})
				return true
			}
			if msg, hit := errorSentinelCompare(pass, bin); hit {
				out = append(out, goidiomatic.Finding{
					RuleID:   r.ID(),
					Message:  msg,
					Pos:      pass.Fset.Position(bin.Pos()),
					Severity: r.Severity(),
				})
			}
			return true
		})
	}
	return out, nil
}

// errorStringCompare matches `x.Error() == "..."` or the mirrored form.
func errorStringCompare(bin *ast.BinaryExpr) (string, bool) {
	if isErrorCallVsStringLit(bin.X, bin.Y) || isErrorCallVsStringLit(bin.Y, bin.X) {
		return "compare errors with errors.Is, not err.Error() against a string literal", true
	}
	return "", false
}

func isErrorCallVsStringLit(a, b ast.Expr) bool {
	call, ok := a.(*ast.CallExpr)
	if !ok || len(call.Args) != 0 {
		return false
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Error" {
		return false
	}
	lit, ok := b.(*ast.BasicLit)
	return ok && lit.Kind == token.STRING
}

// errorSentinelCompare matches `err == someErr` where at least one side has
// type `error` and neither side is the untyped nil. Comparing to nil is
// idiomatic and must not flag.
func errorSentinelCompare(pass *analysis.Pass, bin *ast.BinaryExpr) (string, bool) {
	if isNilIdent(bin.X) || isNilIdent(bin.Y) {
		return "", false
	}
	if !isErrorTyped(pass, bin.X) && !isErrorTyped(pass, bin.Y) {
		return "", false
	}
	return "compare errors with errors.Is, not ==; wrapped errors will not match", true
}

func isNilIdent(e ast.Expr) bool {
	id, ok := e.(*ast.Ident)
	return ok && id.Name == "nil"
}

func isErrorTyped(pass *analysis.Pass, e ast.Expr) bool {
	if pass.TypesInfo == nil {
		return false
	}
	t := pass.TypesInfo.TypeOf(e)
	if t == nil {
		return false
	}
	named, ok := t.(*types.Named)
	if !ok {
		// also accept the predeclared error interface directly
		if iface, ok := t.Underlying().(*types.Interface); ok {
			return iface.NumMethods() == 1 && iface.Method(0).Name() == "Error"
		}
		return false
	}
	if named.Obj() != nil && named.Obj().Name() == "error" {
		return true
	}
	if iface, ok := named.Underlying().(*types.Interface); ok {
		for i := 0; i < iface.NumMethods(); i++ {
			if iface.Method(i).Name() == "Error" {
				return true
			}
		}
	}
	return false
}
