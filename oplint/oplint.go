// Package oplint defines an Analyzer that reports errs.Op const declarations
package oplint

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"strings"
)

var Analyzer = &analysis.Analyzer{
	Name: "oplint",
	Doc:  "reports errs.Op const declarations",
	Run:  run,
}

// render returns the pretty-print of the given node
func render(fset *token.FileSet, x interface{}) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, x); err != nil {
		panic(err)
	}
	return buf.String()
}
func run(pass *analysis.Pass) (interface{}, error) {

	for _, file := range pass.Files {
		allFuncs := collectFuncDeclarations(pass, file)
		for _, f := range allFuncs {
			fo := newFuncOp(f)
			fo.Report(pass)
		}
	}
	return nil, nil
}

// retrieve list of all function declarations
func collectFuncDeclarations(pass *analysis.Pass, node ast.Node) []*ast.FuncDecl {

	// walk through file and pull out func identifier nodes
	var funcDecls []*ast.FuncDecl
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			funcDecls = append(funcDecls, x)
		}
		return true
	})

	return funcDecls
}

type funcOp struct {
	// Function name
	FuncName string
	// Op constant name (either op or Op)
	OpName string
	// Op constant position
	OpNamePos token.Pos
	// Op constant value
	OpValue string
	// Op constant value position
	OpValuePos token.Pos
}

func (fo *funcOp) Report(pass *analysis.Pass) {
	if fo.OpValue != fo.FuncName {
		fmt.Printf("op constant value (%s) at %s does not match the function name (%s)\n", fo.OpValue, pass.Fset.Position(fo.OpNamePos), fo.FuncName)
	}

}

func newFuncOp(fd *ast.FuncDecl) funcOp {

	var fo funcOp

	// get function name
	fo.FuncName = fd.Name.String()

	// loop through list of statements in the function body and
	// determine if any are constant declarations. If so and they
	// have the name op or Op, add to the function and return
	for _, stmt := range fd.Body.List {
		switch x := stmt.(type) {
		case *ast.DeclStmt:
			switch y := x.Decl.(type) {
			case *ast.GenDecl:
				if y.Tok == token.CONST {
					for _, spec := range y.Specs {
						switch z := spec.(type) {
						case *ast.ValueSpec:
							name := z.Names[0]
							if name.Name == "op" || name.Name == "Op" {
								fo.OpName = name.Name
								fo.OpNamePos = name.NamePos
							}
							value := z.Values[0]
							switch lit := value.(type) {
							case *ast.BasicLit:
								fo.OpValue = strings.Trim(lit.Value, "\"")
								fo.OpValuePos = lit.ValuePos
							}
						}
					}
				}
			}
		}
	}

	return fo
}

func nodePrinter(pass *analysis.Pass, node ast.Node) {
	ast.Inspect(node, func(n ast.Node) bool {

		// print type for each ast.Node
		if n != nil {
			fmt.Printf("%T\n", n)
		}

		var (
			s    string
			kind string
		)
		switch x := n.(type) {
		// A BasicLit node represents a literal of basic type.
		case *ast.BasicLit:
			s = x.Value
			kind = "BasicLit"
		// Identifiers
		case *ast.Ident:
			s = x.Name
			if x.Obj != nil {
				kind = x.Obj.Kind.String()
			}
		}
		if s != "" {
			//fmt.Printf("%s:\t%s\n", pass.Fset.Position(n.Pos()), s)
			fmt.Printf("%s[%s]:\t%s\n", pass.Fset.Position(n.Pos()), kind, s)
		}

		return true
	})
}
