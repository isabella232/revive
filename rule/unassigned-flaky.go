package rule

import (
	"fmt"
	"go/ast"

	"github.com/mgechev/revive/lint"
)

// UnassignedFlakyRule lints given else constructs.
type UnassignedFlakyRule struct{}

// Apply applies the rule to given file.
func (r *UnassignedFlakyRule) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	astFile := file.AST
	w := &lintUnassignedFlaky{astFile, onFailure}
	ast.Walk(w, astFile)
	return failures
}

// Name returns the rule name.
func (r *UnassignedFlakyRule) Name() string {
	return "unassigned-flaky"
}

type lintUnassignedFlaky struct {
	file      *ast.File
	onFailure func(lint.Failure)
}

func (w *lintUnassignedFlaky) Visit(node ast.Node) ast.Visitor {
	switch v := node.(type) {
	case *ast.ExprStmt:
		if c, ok := v.X.(*ast.CallExpr); ok {
			if stmt, ok := c.Fun.(*ast.SelectorExpr); ok {
				if stmt.Sel.Name == "Flaky" {
					fmt.Printf("Found statement: %s.%s\n", stmt.X, stmt.Sel)
					w.onFailure(lint.Failure{
						Confidence: .9,
						Category:   "errors",
						Node:       node,
						Failure:    "flake.Flaky(t) needs to reassign the argument -> t := flake.Flaky(t)",
					})
				}
			}
		}
	}
	return w
}
