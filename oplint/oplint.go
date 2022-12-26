// Package oplint defines an Analyzer that reports Op const declarations
// When creating errors or logs, it is common to use an "op" constant to
// capture what function is making the call, to allow for debugging where
// a certain error took place.
package oplint

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "oplint",
	Doc:  "reports inconsistent op const declarations",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		//ast.Print(pass.Fset, file)
		allFuncs := collectFuncDeclarations(file)
		for _, f := range allFuncs {
			fo := newFuncOp(f)
			fo.report(pass)
		}
	}

	return nil, nil
}

// retrieve list of all function declarations
func collectFuncDeclarations(node ast.Node) []*ast.FuncDecl {
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
	// Function Receiver Name
	FuncReceiverIdentifier string
	// Function Receiver Name
	FuncReceiverStructName string
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

func (fo *funcOp) report(pass *analysis.Pass) {

	// only issue warnings for functions that have an op constant defined
	// if no op value, skip
	if fo.OpValue == "" {
		return
	}
	var name string
	name = pass.Pkg.Name() + "/" + fo.FuncName
	if fo.FuncReceiverStructName != "" {
		name = pass.Pkg.Name() + "/" + fo.FuncReceiverStructName + "." + fo.FuncName
	}

	if fo.OpValue != name {
		pass.Reportf(fo.OpNamePos, "%s constant value (%s) does not match function name (%s)", fo.OpName, fo.OpValue, name)
	}
}

func newFuncOp(fd *ast.FuncDecl) funcOp {
	var fo funcOp

	// get function receiver name (if present)
	if fd.Recv != nil {
		field := fd.Recv.List[0]
		fieldName := field.Names[0]
		fo.FuncReceiverIdentifier = fieldName.Name
		switch t := field.Type.(type) {
		case *ast.Ident: // value receiver
			fo.FuncReceiverStructName = t.Name
		case *ast.StarExpr: // pointer receiver
			switch starx := t.X.(type) {
			case *ast.Ident:
				fo.FuncReceiverStructName = starx.Name
			}
		}
	}

	// get function name
	fo.FuncName = fd.Name.String()

	// loop through list of statements in the function body and
	// determine if any are constant declarations. If yes, and they
	// have the name op or Op, add to the struct and return
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
