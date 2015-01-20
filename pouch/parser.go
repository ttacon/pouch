package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strconv"
	"strings"

	"github.com/ttacon/pouch/pouch/defs"
)

func StructsFromFile(pkgName string) ([]*defs.StructInfo, error) {
	fset := token.NewFileSet()
	pf, err := parser.ParseDir(
		fset,
		pkgName,
		func(f os.FileInfo) bool {
			return true
		},
		parser.ParseComments)
	if err != nil {
		return nil, err
	}
	var pkg *ast.Package
	for _, pkg = range pf {
		break
	}

	f := ast.MergePackageFiles(pkg, 0)
	s := &structCollector{}
	ast.Inspect(f, s.Visit)

	return s.structs, nil
}

type structCollector struct {
	structs []*defs.StructInfo
}

func (s *structCollector) Visit(node ast.Node) bool {
	info := structInfo(node)
	if info != nil {
		s.structs = append(s.structs, info)
	}
	return true
}

func structInfo(node ast.Node) *defs.StructInfo {
	// currently I think we'll only be passing in GenDecl,
	// so return warning if cannot find StructType
	genDecl, ok := node.(*ast.GenDecl)
	if !ok {
		return nil
	}

	typeSpec, ok := genDecl.Specs[0].(*ast.TypeSpec)
	if !ok {
		return nil
	}

	structType, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return nil
	}

	return &defs.StructInfo{
		Name:   typeSpec.Name.Name,
		Fields: fromFieldList(structType.Fields),
	}
}

// TODO(ttacon): deal with primary keys appropriately
func fromFieldList(fieldList *ast.FieldList) []defs.FieldInfo {
	var fields []defs.FieldInfo
	for _, field := range fieldList.List {
		isPointer, typ := typeInfo(field.Type)
		for _, name := range field.Names {
			fields = append(fields, defs.FieldInfo{
				Name:      name.Name,
				Column:    columnFromField(name.Name, field.Tag),
				IsPointer: isPointer,
				Type:      typ,
			})
		}
	}
	return fields
}

func columnFromField(name string, t *ast.BasicLit) string {
	if t != nil {
		tag := fromTag(t.Value, "db")
		if len(tag) > 0 {
			return tag
		}
	}

	return name
}

func typeInfo(expr ast.Expr) (bool, string) {
	if id, ok := expr.(*ast.Ident); ok {
		return strings.HasPrefix(id.Name, "*"), strings.TrimPrefix(id.Name, "*")
	}

	if star, ok := expr.(*ast.StarExpr); ok {
		if sel, ok := star.X.(*ast.SelectorExpr); ok {
			pkg, _ := sel.X.(*ast.Ident)
			return true, pkg.Name + "." + sel.Sel.Name
		}
		if id, ok := star.X.(*ast.Ident); ok {
			return true, id.Name
		}
	}

	return false, ""
}

// taken from http://golang.org/src/reflect/type.go?s=21589:21632#L752
func fromTag(tag, key string) string {
	tag = tag[1 : len(tag)-1]
	for tag != "" {
		// skip leading space
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// scan to colon.
		// a space or a quote is a syntax error
		i = 0
		for i < len(tag) && tag[i] != ' ' && tag[i] != ':' && tag[i] != '"' {
			i++
		}
		if i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// scan quoted string to find value
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		if key == name {
			value, _ := strconv.Unquote(qvalue)
			return value
		}
	}
	return ""
}
