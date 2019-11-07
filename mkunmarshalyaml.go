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
	"strconv"
	"strings"
)

func main() {
	var buf bytes.Buffer
	outf := func(format string, args ...interface{}) {
		fmt.Fprintf(&buf, format, args...)
	}

	outf("// Code generated by mkunmarshalyaml.go. DO NOT EDIT.")
	outf("\n\npackage openapi")
	outf("\n\nimport (")
	outf("\n\"errors\"")
	outf("\n\"net/url\"")
	outf("\n\"regexp\"")
	outf("\n\"strconv\"")
	outf("\n\"strings\"")
	outf("\n")
	outf("\nyaml \"github.com/goccy/go-yaml\"")
	outf("\n)")

	f, err := parser.ParseFile(token.NewFileSet(), "interfaces.go", nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if genDecl.Doc == nil || len(genDecl.Doc.List) == 0 || genDecl.Doc.List[0].Text != "//+object" {
			log.Printf("%v is not an openapi object. skip.", genDecl.Specs[0].(*ast.TypeSpec).Name.Name)
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

			log.Printf("generate %s.Unmarshal", typ.Name.Name)
			outf("\n\nfunc (v *%s) UnmarshalYAML(b []byte) error {", typ.Name.Name)
			outf("\nvar proxy map[string]raw")
			outf("\nif err := yaml.Unmarshal(b, &proxy); err != nil {")
			outf("\nreturn err")
			outf("\n}")

			for _, field := range st.Fields.List {
				fn := field.Names[0].Name
				tag := parseTags(field)
				yn := yamlName(field, tag)
				required := isRequired(tag)

				if yn == "-" {
					continue
				}

				if isInline(tag) {
					ft, ok := field.Type.(*ast.MapType)
					if !ok {
						log.Fatalf("exected map for inline %s but %s", yn, field.Type)
					}
					formatTag := tag["format"]
					outf("\n%s := map[string]%s{}", fn, ast2type(ft.Value))
					outf("\nfor key, val := range proxy {")
					if len(formatTag) > 0 {
						switch formatTag[0] {
						case "prefix":
							outf("\nif !strings.HasPrefix(key, \"%s\") {", formatTag[1])
							outf("\ncontinue")
							outf("\n}")
						case "regexp":
							outf("\n%sRegexp := regexp.MustCompile(`%s`)", fn, formatTag[1])
							outf("\nif !%sRegexp.MatchString(key) {", fn)
							outf("\ncontinue")
							outf("\n}")
						}
					}
					outf("\nvar %sv %s", fn, strings.TrimPrefix(ast2type(ft.Value), "*"))
					outf("\nif err := yaml.Unmarshal(val, &%sv); err != nil {", fn)
					outf("\nreturn err")
					outf("\n}")
					outf("\n%s[key] = ", fn)
					if _, ok := ft.Value.(*ast.StarExpr); ok {
						outf("&")
					}
					outf("%sv", fn)
					outf("\n}")
					outf("\nif len(%s) != 0 {", fn)
					outf("\nv.%s = %s", fn, fn)
					outf("\n}")
					continue
				}

				unmarshalField := func() {
					outf("\nif err := yaml.Unmarshal(%sBytes, &%[1]s); err != nil {", fn)
					outf("\nreturn err")
					outf("\n}")
				}

				outf("\n\n")
				if required {
					outf("%sBytes, ok := proxy[\"%s\"]", fn, yn)
					outf("\nif !ok {")
					outf("\nreturn errors.New(`\"%s\" field is required`)", yn)
					outf("\n}")
				} else {
					outf("if %sBytes, ok := proxy[\"%s\"]; ok {", fn, yn)
				}

				switch t := field.Type.(type) {
				case *ast.Ident: // built-in type
					switch t.Name {
					case "string":
						outf("\nv.%s = string(%[1]sBytes)", fn)
					case "bool":
						outf("\nt, err := strconv.ParseBool(string(%sBytes))", fn)
						outf("\nif err != nil {")
						outf("\nreturn err")
						outf("\n}")
						outf("\nv.%s = t", fn)
					case "int":
						outf("\ni, err := strconv.Atoi(string(%sBytes))", fn)
						outf("\nif err != nil {")
						outf("\nreturn err")
						outf("\n}")
						outf("\nv.%s = i", fn)
					default:
						log.Fatalf("unknown type for %s: %s", fn, t.Name)
					}
				default:
					outf("\nvar %s %s", fn, strings.TrimPrefix(ast2type(t), "*"))
					unmarshalField()
					outf("\nv.%s = ", fn)
					if _, ok := t.(*ast.StarExpr); ok {
						outf("&")
					}
					outf("%s", fn)
				}
				if !required {
					outf("\n}")
				}
				formatValidation(outf, fn, yn, field, tag, required)
			}

			outf("\nreturn nil")
			outf("\n}")
		}
	}
	src, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatalf("error on formatting: %+v", err)
	}
	if err := ioutil.WriteFile("unmarshalyaml_gen.go", src, 0644); err != nil {
		log.Fatal(err)
	}
}

func formatValidation(outf func(string, ...interface{}), fieldname, yamlname string, field *ast.Field, tag tags, required bool) {
	switch tag.get("format") {
	case "semver":
		outf("\n\nif !isValidSemVer(v.%s) {", fieldname)
		outf("\nreturn errors.New(`\"%s\" field must be a valid semantic version but not`)", yamlname)
		outf("\n}")
	case "url":
		outf("\n")
		if !required {
			outf("\nif v.%s != \"\" {", fieldname)
		}
		outf("\nif _, err := ")
		if len(tag["format"]) > 1 && tag["format"][1] == "template" {
			outf("url.Parse(urlTemplateVarRegexp.ReplaceAllLiteralString(v.%s, `placeholder`))", fieldname)
		} else {
			outf("url.ParseRequestURI(v.%s)", fieldname)
		}
		outf("; err != nil {")
		outf("\nreturn err")
		outf("\n}")
		if !required {
			outf("\n}")
		}
	case "email":
		outf("\n")
		if !required {
			outf("\nif v.%s != \"\" {", fieldname)
		}
		outf("\n\nif v.%s != \"\" && !emailRegexp.MatchString(v.%[1]s) {", fieldname)
		outf("\nreturn errors.New(`\"%s\" field must be an email address`)", yamlname)
		outf("\n}")
		if !required {
			outf("\n}")
		}
	case "runtime":
		if _, ok := field.Type.(*ast.MapType); ok {
			outf("\n\nfor key := range v.%s {", fieldname)
			outf("\nif !matchRuntimeExpr(key) {")
			outf("\nreturn errors.New(`the keys of \"%s\" must be a runtime expression`)", yamlname)
			outf("\n}")
			outf("\n}")
		}
	case "regexp":
		if _, ok := field.Type.(*ast.MapType); ok {
			outf("\n\n%sRegexp := regexp.MustCompile(`%s`)", fieldname, tag["format"][1])
			outf("\nfor key := range v.%s {", fieldname)
			outf("\nif !%sRegexp.MatchString(v.%s) {", fieldname, fieldname)
			outf("\nreturn errors.New(`the keys of \"%s\" must be match \"%s\"`)", yamlname, tag["format"][1])
			outf("\n}")
		}
	}
	if list, ok := tag["oneof"]; ok {
		outf("\n")
		if !required {
			outf("\nif v.%s != \"\" {", fieldname)
		}
		outf("\nif isOneOf(v.%s, %#v) {", fieldname, list)
		outf("\nreturn errors.New(`\"%s\" field must be one of [%s]`)", yamlname, strings.Join(quoteEachString(list), ", "))
		outf("\n}")
		if !required {
			outf("\n}")
		}
	}
}

func isInline(t tags) bool {
	vs := t["yaml"]
	if len(vs) < 2 {
		return false
	}
	for _, v := range vs[1:] {
		if v == "inline" {
			return true
		}
	}
	return false
}

func isRequired(t tags) bool {
	return t.get("required") != ""
}

func yamlName(field *ast.Field, t tags) string {
	yn := t.get("yaml")
	if yn != "" {
		return yn
	}
	return field.Names[0].Name
}

func quoteEachString(list []string) []string {
	ret := make([]string, len(list))
	for i := range list {
		ret[i] = strconv.Quote(list[i])
	}
	return ret
}