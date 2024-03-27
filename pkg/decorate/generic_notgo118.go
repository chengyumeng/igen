//go:build !go1.18
// +build !go1.18

package decorate

import (
	"go/ast"

	"github.com/chengyumeng/igen/pkg/decorate/model"
)

func getTypeSpecTypeParams(ts *ast.TypeSpec) []*ast.Field {
	return nil
}

func (p *fileParser) parseGenericType(pkg string, typ ast.Expr, tps map[string]bool) (model.Type, error) {
	return nil, nil
}

func getIdentTypeParams(decl interface{}) string {
	return ""
}
