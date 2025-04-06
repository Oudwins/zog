package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log/slog"
	"os"
	"strconv"

	"github.com/Oudwins/zog/cmd/pkg/assert"
	"github.com/Oudwins/zog/cmd/pkg/logger"
)

func getExecCommands(src []byte, goline int) (string, token.Pos) {
	// Directly scan the source file for the go:generate directive
	lines := string(src)
	lineStart := 0
	currentLine := 1

	// var execCommands string
	var generatePos token.Pos
	var generateCommands string

	for i, char := range lines {
		if char == '\n' {
			// Check if this is the line we're looking for
			if currentLine == goline {
				lineContent := lines[lineStart:i]
				slog.Debug("Found target line", "line", goline, "content", lineContent)

				// Check if this line contains a go:generate directive
				if len(lineContent) >= 13 && lineContent[:13] == "//go:generate" {
					generatePos = token.Pos(i + 1) // Position after the directive
					generateCommands = lineContent[13:]
					break
				}
			}

			lineStart = i + 1
			currentLine++
		}
	}

	return generateCommands, generatePos
}

func main() {
	logger.Init()
	// Retrieve the filename and line number from environment variables
	gofile := os.Getenv("GOFILE")
	goline, err := strconv.Atoi(os.Getenv("GOLINE"))
	assert.NotError(err, "Error converting GOLINE to integer")

	// Open the source file
	src, err := os.ReadFile(gofile)
	assert.NotError(err, "Error reading source file")
	slog.Debug("Read source file", "file", string(src))

	generateCommands, generatePos := getExecCommands(src, goline)
	slog.Debug("Generate commands", "commands", generateCommands, "pos", generatePos)
	//

	// Create the file set and parse the file
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, gofile, src, parser.ParseComments)
	assert.NotError(err, "Error parsing source file")

	// Find the first declaration after the generate position
	var nextDecl ast.Decl
	for _, decl := range file.Decls {
		if decl.Pos() > generatePos {
			nextDecl = decl
			break
		}
	}

	if nextDecl == nil {
		slog.Debug("No declaration found after //go:generate directive")
		return
	}

	slog.Debug("Found next declaration", "content", string(src[nextDecl.Pos()-1:nextDecl.End()]))

	ast.Inspect(nextDecl, func(n ast.Node) bool {
		if n == nil {
			return false
		}
		slog.Debug("Found node", "value", fmt.Sprintf("%+v", n))
		return true
	})

	// printer.Fprint(writer, fset, nextDecl)

	// Log the type and details of the declaration
	// var logAstNode func(node ast.Node, depth int)
	// logAstNode = func(node ast.Node, depth int) {
	// 	if node == nil {
	// 		return
	// 	}

	// 	txt := src[node.Pos():node.End()]

	// 	slog.Debug("Found Node",
	// 		"value", fmt.Sprintf("%+v", node),
	// 		"text", txt,
	// 	)

	// 	switch n := node.(type) {
	// 	case *ast.GenDecl:
	// 		slog.Debug("Found GenDecl",
	// 			"kind", n.Tok.String(),
	// 			"specs_count", len(n.Specs),
	// 			"value", fmt.Sprintf("%+v", n),
	// 		)

	// 		for _, spec := range n.Specs {
	// 			logAstNode(spec, depth+1)
	// 		}

	// 	case *ast.ValueSpec:
	// 		slog.Debug("ValueSpec",
	// 			"names", fmt.Sprintf("%+v", n.Names),
	// 			"type", fmt.Sprintf("%+v", n.Type),
	// 			"values", fmt.Sprintf("%+v", n.Values),
	// 			"comments", fmt.Sprintf("%+v", n.Comment))

	// 		logAstNode(n.Type, depth+1)
	// 		for _, value := range n.Values {
	// 			logAstNode(value, depth+1)
	// 		}

	// 	case *ast.TypeSpec:
	// 		slog.Debug("TypeSpec",
	// 			"name", n.Name.Name,
	// 			"type", fmt.Sprintf("%#v", n.Type),
	// 			"comments", fmt.Sprintf("%#v", n.Comment))

	// 		logAstNode(n.Type, depth+1)

	// 	case *ast.FuncDecl:
	// 		slog.Debug("FuncDecl",
	// 			"name", n.Name.Name,
	// 			"recv", fmt.Sprintf("%#v", n.Recv),
	// 			"type", fmt.Sprintf("%#v", n.Type),
	// 			"body", fmt.Sprintf("%#v", n.Body))

	// 		logAstNode(n.Recv, depth+1)
	// 		logAstNode(n.Type, depth+1)
	// 		logAstNode(n.Body, depth+1)

	// 	case *ast.FuncType:
	// 		slog.Debug("FuncType",
	// 			"params", fmt.Sprintf("%#v", n.Params),
	// 			"results", fmt.Sprintf("%#v", n.Results))

	// 		logAstNode(n.Params, depth+1)
	// 		logAstNode(n.Results, depth+1)

	// 	case *ast.FieldList:
	// 		if n != nil {
	// 			slog.Debug("FieldList",
	// 				"fields", fmt.Sprintf("%#v", n.List))
	// 			for _, field := range n.List {
	// 				logAstNode(field, depth+1)
	// 			}
	// 		}

	// 	case *ast.Field:
	// 		slog.Debug("Field",
	// 			"names", fmt.Sprintf("%#v", n.Names),
	// 			"type", fmt.Sprintf("%#v", n.Type),
	// 			"tag", fmt.Sprintf("%#v", n.Tag))

	// 		logAstNode(n.Type, depth+1)

	// 	case *ast.BlockStmt:
	// 		slog.Debug("BlockStmt",
	// 			"statements", fmt.Sprintf("%#v", n.List))
	// 		for _, stmt := range n.List {
	// 			logAstNode(stmt, depth+1)
	// 		}

	// 	case *ast.BasicLit:
	// 		slog.Debug("BasicLit",
	// 			"kind", n.Kind.String(),
	// 			"value", n.Value)

	// 	case *ast.Ident:
	// 		slog.Debug("Ident",
	// 			"name", n.Name,
	// 			"obj", fmt.Sprintf("%#v", n.Obj))

	// 	case *ast.CallExpr:
	// 		slog.Debug("Other Node",
	// 			"type", fmt.Sprintf("%+T", n),
	// 			"value", fmt.Sprintf("%+v", n))

	// 	default:
	// 		slog.Debug("Other Node",
	// 			"type", fmt.Sprintf("%+T", n),
	// 			"value", fmt.Sprintf("%+v", n))

	// 	}
	// }

	// logAstNode(nextDecl, 0)

	// Find the declaration immediately following the //go:generate directive
	// var targetDecl ast.Decl
	// for _, decl := range file.Decls {
	// 	if decl.Pos() > generatePos {
	// 		targetDecl = decl
	// 		break
	// 	}
	// }

	// if targetDecl == nil {
	// 	fmt.Println("No declaration found after //go:generate directive.")
	// 	os.Exit(1)
	// }

	// // Process the found declaration (e.g., print its type)
	// switch decl := targetDecl.(type) {
	// case *ast.GenDecl:
	// 	for _, spec := range decl.Specs {
	// 		if typeSpec, ok := spec.(*ast.TypeSpec); ok {
	// 			fmt.Printf("Found type declaration: %s\n", typeSpec.Name.Name)
	// 			// Further processing can be done here
	// 		}
	// 	}
	// default:
	// 	fmt.Println("Declaration is not a type declaration.")
	// }
}
