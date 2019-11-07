// +build ignore

package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"strings"
	"unicode"
)

func main() {
	var buf bytes.Buffer
	outf := func(format string, args ...interface{}) {
		fmt.Fprintf(&buf, format, args...)
	}

	outf("// Code generated by mkgetter.go. DO NOT EDIT.")
	outf("\n\npackage openapi")

	f, err := parser.ParseFile(token.NewFileSet(), "interfaces.go", nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range genDecl.Specs {
			typ, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			st, ok := typ.Type.(*ast.StructType)
			if !ok {
				continue
			}

			for _, field := range st.Fields.List {
				fn := field.Names[0].Name
				ft := ast2type(field.Type)

				outf("\n\nfunc (v *%s) %s() %s {", typ.Name.Name, expose(fn), ft)

				if strings.HasPrefix(ft, "*") {
					outf("\nif v.%s == nil {", fn)
					// trim "*"
					outf("\nreturn &%s{}", ft[1:])
					outf("\n}")
				}

				outf("\nreturn v.%s", fn)
				outf("\n}")
			}
		}
	}
	src, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile("getter_gen.go", src, 0644); err != nil {
		log.Fatal(err)
	}
}

func expose(ident string) string {
	rident := []rune(ident)
	return string(append([]rune{unicode.ToUpper(rident[0])}, rident[1:]...))
}