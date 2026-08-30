package codegen

import (
	"fmt"
	"go/ast"
	"go/types"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

// ZogSchemas returns the concrete schema receiver types declared in dir.
func ZogSchemas(dir string) ([]string, error) {
	loaded, err := packages.Load(&packages.Config{Mode: packages.LoadSyntax, Dir: dir}, ".")
	if err != nil {
		return nil, fmt.Errorf("load package: %w", err)
	}
	if len(loaded) != 1 || loaded[0].Name != "zog" {
		return nil, fmt.Errorf("package zog not found in %s", dir)
	}
	pkg := loaded[0]
	if len(pkg.Errors) > 0 {
		return nil, fmt.Errorf("load package: %s", pkg.Errors[0])
	}

	structs := map[string][]string{}

	for _, file := range pkg.Syntax {
		for _, declaration := range file.Decls {
			declaration, ok := declaration.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, spec := range declaration.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
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
				structs[typeSpec.Name.Name] = typeParams
			}
		}
	}

	zogSchemaObject := pkg.Types.Scope().Lookup("ZogSchema")
	if zogSchemaObject == nil {
		return nil, fmt.Errorf("ZogSchema not found")
	}
	zogSchema, ok := zogSchemaObject.Type().Underlying().(*types.Interface)
	if !ok {
		return nil, fmt.Errorf("ZogSchema is not an interface")
	}

	var schemas []string
	for typeName, typeParams := range structs {
		named, ok := pkg.Types.Scope().Lookup(typeName).Type().(*types.Named)
		if !ok || !types.Implements(types.NewPointer(named), zogSchema) {
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
