package rules

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"

	goidiomatic "github.com/CivNode/go-idiomatic"
)

// ContextFirstArg flags functions that take context.Context as any
// parameter other than the first. The Go convention is that Context is
// always the first parameter, named ctx.
var ContextFirstArg goidiomatic.Rule = contextFirstArg{}

type contextFirstArg struct{}

func (contextFirstArg) ID() string   { return "context-first-arg" }
func (contextFirstArg) Name() string { return "context.Context must be first parameter" }
func (contextFirstArg) Description() string {
	return "By convention context.Context is always the first parameter of a function, conventionally named ctx."
}
func (contextFirstArg) Severity() goidiomatic.Severity { return goidiomatic.Warn }

func (r contextFirstArg) Check(pass *analysis.Pass) ([]goidiomatic.Finding, error) {
	var out []goidiomatic.Finding
	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			var params *ast.FieldList
			switch v := n.(type) {
			case *ast.FuncDecl:
				if v.Type != nil {
					params = v.Type.Params
				}
			case *ast.FuncLit:
				if v.Type != nil {
					params = v.Type.Params
				}
			default:
				return true
			}
			if params == nil || len(params.List) < 2 {
				return true
			}

			// Position i counts individual parameters, not Field groups, so
			// `func(a, b int, c int)` reports indices 0,1,2 correctly.
			idx := 0
			for _, field := range params.List {
				n := len(field.Names)
				if n == 0 {
					n = 1
				}
				if idx > 0 && isContextType(pass, field.Type) {
					out = append(out, goidiomatic.Finding{
						RuleID:   r.ID(),
						Message:  "context.Context must be the first parameter",
						Pos:      pass.Fset.Position(field.Pos()),
						Severity: r.Severity(),
					})
					break
				}
				idx += n
			}
			return true
		})
	}
	return out, nil
}

func isContextType(pass *analysis.Pass, expr ast.Expr) bool {
	// Fast syntactic path: `context.Context`.
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
	obj := named.Obj()
	if obj == nil || obj.Name() != "Context" {
		return false
	}
	if obj.Pkg() == nil {
		return false
	}
	return obj.Pkg().Path() == "context"
}
