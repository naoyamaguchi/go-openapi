// +build ignore

package main

import (
	"go/ast"
	"log"
	"reflect"
)

func ast2type(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr: // pointer
		return "*" + t.X.(*ast.Ident).Name
	case *ast.ArrayType:
		return "[]" + ast2type(t.Elt)
	case *ast.MapType:
		return "map[" + ast2type(t.Key) + "]" + ast2type(t.Value)
	case *ast.InterfaceType:
		return "interface{}"
	default:
		log.Fatalf("unknown type: %s", reflect.TypeOf(t))
	}
	return ""
}
