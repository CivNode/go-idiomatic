package rules

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"

	goidiomatic "github.com/CivNode/go-idiomatic"
)

// NoSleepForCoordination flags time.Sleep calls inside a function that
// also uses a channel operation or a context.Context. Sleeping as a way
// to wait for "the other thing" to happen is almost always a race.
var NoSleepForCoordination goidiomatic.Rule = noSleepForCoordination{}

type noSleepForCoordination struct{}

func (noSleepForCoordination) ID() string   { return "no-sleep-for-coordination" }
func (noSleepForCoordination) Name() string { return "no time.Sleep for coordination" }
func (noSleepForCoordination) Description() string {
	return "Using time.Sleep to wait for another goroutine or a context deadline is a race. Prefer channel receives, sync primitives, or context.Done."
}
func (noSleepForCoordination) Severity() goidiomatic.Severity { return goidiomatic.Error }

func (r noSleepForCoordination) Check(pass *analysis.Pass) ([]goidiomatic.Finding, error) {
	var out []goidiomatic.Finding
	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			var body *ast.BlockStmt
			switch v := n.(type) {
			case *ast.FuncDecl:
				body = v.Body
				if body != nil && v.Type != nil {
					if hasContextParam(pass, v.Type.Params) {
						if pos, ok := findSleep(body); ok {
							out = append(out, goidiomatic.Finding{
								RuleID:   r.ID(),
								Message:  "time.Sleep alongside context.Context looks like coordination; use ctx.Done()",
								Pos:      pass.Fset.Position(pos),
								Severity: r.Severity(),
							})
							return true
						}
					}
				}
			case *ast.FuncLit:
				body = v.Body
				if body != nil && v.Type != nil && hasContextParam(pass, v.Type.Params) {
					if pos, ok := findSleep(body); ok {
						out = append(out, goidiomatic.Finding{
							RuleID:   r.ID(),
							Message:  "time.Sleep alongside context.Context looks like coordination; use ctx.Done()",
							Pos:      pass.Fset.Position(pos),
							Severity: r.Severity(),
						})
						return true
					}
				}
			default:
				return true
			}

			if body == nil {
				return true
			}
			if !usesChannel(body) {
				return true
			}
			if pos, ok := findSleep(body); ok {
				out = append(out, goidiomatic.Finding{
					RuleID:   r.ID(),
					Message:  "time.Sleep alongside channel operations looks like coordination; use a channel receive or select",
					Pos:      pass.Fset.Position(pos),
					Severity: r.Severity(),
				})
			}
			return true
		})
	}
	return out, nil
}

func findSleep(body *ast.BlockStmt) (token.Pos, bool) {
	var pos token.Pos
	var ok bool
	ast.Inspect(body, func(n ast.Node) bool {
		call, isCall := n.(*ast.CallExpr)
		if !isCall {
			return true
		}
		sel, isSel := call.Fun.(*ast.SelectorExpr)
		if !isSel {
			return true
		}
		id, isID := sel.X.(*ast.Ident)
		if !isID {
			return true
		}
		if id.Name == "time" && sel.Sel.Name == "Sleep" {
			pos = call.Pos()
			ok = true
			return false
		}
		return true
	})
	return pos, ok
}

func usesChannel(body *ast.BlockStmt) bool {
	var hit bool
	ast.Inspect(body, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.SendStmt, *ast.SelectStmt, *ast.ChanType:
			_ = v
			hit = true
			return false
		case *ast.UnaryExpr:
			if v.Op == token.ARROW {
				hit = true
				return false
			}
		}
		return true
	})
	return hit
}

func hasContextParam(pass *analysis.Pass, params *ast.FieldList) bool {
	if params == nil {
		return false
	}
	for _, field := range params.List {
		if isContextParamType(pass, field.Type) {
			return true
		}
	}
	return false
}

func isContextParamType(pass *analysis.Pass, expr ast.Expr) bool {
	if sel, ok := expr.(*ast.SelectorExpr); ok {
		if id, ok := sel.X.(*ast.Ident); ok {
			if id.Name == "context" && sel.Sel.Name == "Context" {
				return true
			}
		}
	}
	if pass.TypesInfo == nil {
		return false
	}
	t := pass.TypesInfo.TypeOf(expr)
	if t == nil {
		return false
	}
	named, ok := t.(*types.Named)
	if !ok {
		return false
	}
	if named.Obj() == nil || named.Obj().Name() != "Context" {
		return false
	}
	if named.Obj().Pkg() == nil {
		return false
	}
	return named.Obj().Pkg().Path() == "context"
}
