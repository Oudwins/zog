package codegen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"sort"
	"strings"
)

// ZogSchemas returns the concrete schema receiver types declared in dir.
func ZogSchemas(dir string) ([]string, error) {
	packages, err := parser.ParseDir(token.NewFileSet(), dir, func(info fs.FileInfo) bool {
		return !strings.HasSuffix(info.Name(), "_test.go")
	}, 0)
	if err != nil {
		return nil, fmt.Errorf("parse package: %w", err)
	}

	pkg, ok := packages["zog"]
	if !ok {
		return nil, fmt.Errorf("package zog not found in %s", dir)
	}

	requiredMethods := map[string]bool{}
	types := map[string][]string{}
	methods := map[string]map[string]bool{}

	for _, file := range pkg.Files {
		for _, declaration := range file.Decls {
			switch declaration := declaration.(type) {
			case *ast.GenDecl:
				for _, spec := range declaration.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					if typeSpec.Name.Name == "ZogSchema" {
						interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
						if !ok {
							return nil, fmt.Errorf("ZogSchema is not an interface")
						}
						for _, method := range interfaceType.Methods.List {
							for _, name := range method.Names {
								requiredMethods[name.Name] = true
							}
						}
					}
					if _, ok := typeSpec.Type.(*ast.StructType); !ok {
						continue
					}
					var typeParams []string
					if typeSpec.TypeParams != nil {
						for _, field := range typeSpec.TypeParams.List {
							for _, name := range field.Names {
								typeParams = append(typeParams, name.Name)
							}
						}
					}
					types[typeSpec.Name.Name] = typeParams
				}
			case *ast.FuncDecl:
				if declaration.Recv == nil || len(declaration.Recv.List) == 0 {
					continue
				}
				typeName := receiverTypeName(declaration.Recv.List[0].Type)
				if typeName == "" {
					continue
				}
				if methods[typeName] == nil {
					methods[typeName] = map[string]bool{}
				}
				methods[typeName][declaration.Name.Name] = true
			}
		}
	}

	if len(requiredMethods) == 0 {
		return nil, fmt.Errorf("ZogSchema has no methods")
	}

	var schemas []string
	for typeName, typeParams := range types {
		if !hasAllMethods(methods[typeName], requiredMethods) {
			continue
		}
		if len(typeParams) > 0 {
			typeName += "[" + strings.Join(typeParams, ", ") + "]"
		}
		schemas = append(schemas, typeName)
	}
	sort.Strings(schemas)
	return schemas, nil
}

func receiverTypeName(expr ast.Expr) string {
	switch expr := expr.(type) {
	case *ast.Ident:
		return expr.Name
	case *ast.StarExpr:
		return receiverTypeName(expr.X)
	case *ast.IndexExpr:
		return receiverTypeName(expr.X)
	case *ast.IndexListExpr:
		return receiverTypeName(expr.X)
	default:
		return ""
	}
}

func hasAllMethods(methods, requiredMethods map[string]bool) bool {
	for method := range requiredMethods {
		if !methods[method] {
			return false
		}
	}
	return true
}
