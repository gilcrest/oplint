// Package oplint defines an Analyzer that reports Op const declarations
// When creating errors or logs, it is common to use an "op" constant to
// capture what function is making the call, to allow for debugging where
// a certain error took place.
package oplint

import (
	"flag"
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func newFlagSet() flag.FlagSet {
	fs := flag.NewFlagSet("oplint flags", flag.ContinueOnError)
	fs.Bool("missing", false, "provide diagnostics for functions which have an error return, but no op constant defined")

	return *fs
}

var Analyzer = &analysis.Analyzer{
	Name:  "oplint",
	Doc:   "reports inconsistent op const declarations",
	Run:   run,
	Flags: newFlagSet(),
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		//ast.Print(pass.Fset, file)
		allFuncs := collectFuncDeclarations(file)
		for _, f := range allFuncs {
			fo := newFuncOp(f)
			fo.reportConstantMismatch(pass)

			// retrieve *flag.Flag from Analyzer flag.Flagset
			mf := pass.Analyzer.Flags.Lookup("missing")
			if mf != nil {
				// perform type assertion to retrieve boolean value
				m := mf.Value.(flag.Getter).Get().(bool)
				if m {
					fo.reportMissingConstant(pass)
				}
			}
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
	// Function Declaration
	FuncDecl *ast.FuncDecl
	// HasErrorResult reports whether function returns an error
	HasErrorResult bool
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

func (fo *funcOp) reportConstantMismatch(pass *analysis.Pass) {

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

	// report op or Op constants defined for functions,
	// but the value does not match the function name
	if fo.OpValue != name {
		pass.Reportf(fo.OpNamePos, "%s constant value (%s) does not match function name (%s)", fo.OpName, fo.OpValue, name)
	}
}

func (fo *funcOp) reportMissingConstant(pass *analysis.Pass) {

	// only issue warnings for functions that have no op constant defined
	// if it has op value, skip
	if fo.OpValue != "" {
		return
	}
	var name string
	name = pass.Pkg.Name() + "/" + fo.FuncName
	if fo.FuncReceiverStructName != "" {
		name = pass.Pkg.Name() + "/" + fo.FuncReceiverStructName + "." + fo.FuncName
	}

	if fo.HasErrorResult && fo.OpValue == "" {
		pass.Reportf(fo.FuncDecl.Pos(), "%s returns an error but does not define an op constant", name)
	}
}

func newFuncOp(fd *ast.FuncDecl) *funcOp {
	fo := &funcOp{FuncDecl: fd}

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

	// determine if function returns an error result
	if fd.Type.Results.NumFields() > 0 {
		for _, rf := range fd.Type.Results.List {
			switch t := rf.Type.(type) {
			case *ast.Ident:
				if t.Name == "error" {
					fo.HasErrorResult = true
				}
			}
		}
	}

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
							// if constant is named op or Op,
							// add name and position to struct
							name := z.Names[0]
							if name.Name == "op" || name.Name == "Op" {
								fo.OpName = name.Name
								fo.OpNamePos = name.NamePos
							}
							if fo.OpName != "" {
								value := z.Values[0]
								switch lit := value.(type) {
								case *ast.BasicLit:
									// add constant op (or Op) value and position to struct
									fo.OpValue = strings.Trim(lit.Value, "\"")
									fo.OpValuePos = lit.ValuePos
								}
							}
						}
					}
				}
			}
		}
	}

	return fo
}
