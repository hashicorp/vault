// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains the pieces of the tool that use typechecking from the go/types package.

package govet

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/importer"
	"go/token"
	"go/types"
	"os"
	"time"

	"golang.org/x/tools/go/gcexportdata"
	"golang.org/x/tools/go/loader"
)

// stdImporter is the importer we use to import packages.
// It is shared so that all packages are imported by the same importer.
var stdImporter types.Importer

var (
	errorType     *types.Interface
	stringerType  *types.Interface // possibly nil
	formatterType *types.Interface // possibly nil
)

func newGCExportDataImporter(fset *token.FileSet) types.ImporterFrom {
	return gcexportdata.NewImporter(fset, make(map[string]*types.Package))
}

func inittypes() error {
	errorType = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)

	typ, err := importType("fmt", "Stringer")
	if err != nil {
		return err
	}
	stringerType = typ.Underlying().(*types.Interface)

	typ, err = importType("fmt", "Formatter")
	if err != nil {
		return err
	}
	formatterType = typ.Underlying().(*types.Interface)

	return nil
}

// isNamedType reports whether t is the named type path.name.
func isNamedType(t types.Type, path, name string) bool {
	n, ok := t.(*types.Named)
	if !ok {
		return false
	}
	obj := n.Obj()
	return obj.Name() == name && obj.Pkg() != nil && obj.Pkg().Path() == path
}

// importType returns the type denoted by the qualified identifier
// path.name, and adds the respective package to the imports map
// as a side effect. In case of an error, importType returns nil.
func importType(path, name string) (types.Type, error) {
	startedAt := time.Now()
	defer func() {
		fmt.Fprintf(os.Stderr, "vet: import of type %s.%s took %s\n", path, name, time.Since(startedAt))
	}()

	pkg, err := stdImporter.Import(path)
	if err != nil {
		// This can happen if the package at path hasn't been compiled yet.
		return nil, fmt.Errorf("import of type %s.%s failed: %v", path, name, err)
	}
	if obj, ok := pkg.Scope().Lookup(name).(*types.TypeName); ok {
		return obj.Type(), nil
	}

	return nil, fmt.Errorf("can't import type %s.%s: invalid type name %q", path, name, name)
}

func (pkg *Package) check(fs *token.FileSet, astFiles []*ast.File, pkgInfo *loader.PackageInfo) []error {
	if stdImporter == nil {
		if *source {
			stdImporter = importer.For("source", nil)
		} else {
			stdImporter = newGCExportDataImporter(fs)
		}
		if err := inittypes(); err != nil {
			return []error{fmt.Errorf("can't init std types: %s", err)}
		}
	}

	var allErrors []error
	pkg.spans = make(map[types.Object]Span)
	if pkgInfo != nil {
		pkg.defs = pkgInfo.Defs
		pkg.uses = pkgInfo.Uses
		pkg.selectors = pkgInfo.Selections
		pkg.types = pkgInfo.Types
		pkg.typesPkg = pkgInfo.Pkg
	} else {
		pkg.defs = make(map[*ast.Ident]types.Object)
		pkg.uses = make(map[*ast.Ident]types.Object)
		pkg.selectors = make(map[*ast.SelectorExpr]*types.Selection)
		pkg.types = make(map[ast.Expr]types.TypeAndValue)

		config := types.Config{
			// We use the same importer for all imports to ensure that
			// everybody sees identical packages for the given paths.
			Importer: stdImporter,
			// By providing a Config with our own error function, it will continue
			// past the first error. We collect them all for printing later.
			Error: func(e error) {
				allErrors = append(allErrors, e)
			},

			Sizes: archSizes,
		}
		info := &types.Info{
			Selections: pkg.selectors,
			Types:      pkg.types,
			Defs:       pkg.defs,
			Uses:       pkg.uses,
		}
		typesPkg, err := config.Check(pkg.path, fs, astFiles, info)
		if len(allErrors) == 0 && err != nil {
			allErrors = append(allErrors, fmt.Errorf("type checker failed: %s", err))
		}

		pkg.typesPkg = typesPkg
	}

	// update spans
	for id, obj := range pkg.defs {
		pkg.growSpan(id, obj)
	}
	for id, obj := range pkg.uses {
		pkg.growSpan(id, obj)
	}
	return allErrors
}

// matchArgType reports an error if printf verb t is not appropriate
// for operand arg.
//
// typ is used only for recursive calls; external callers must supply nil.
//
// (Recursion arises from the compound types {map,chan,slice} which
// may be printed with %d etc. if that is appropriate for their element
// types.)
func (f *File) matchArgType(t printfArgType, typ types.Type, arg ast.Expr) bool {
	return f.matchArgTypeInternal(t, typ, arg, make(map[types.Type]bool))
}

// matchArgTypeInternal is the internal version of matchArgType. It carries a map
// remembering what types are in progress so we don't recur when faced with recursive
// types or mutually recursive types.
func (f *File) matchArgTypeInternal(t printfArgType, typ types.Type, arg ast.Expr, inProgress map[types.Type]bool) bool {
	// %v, %T accept any argument type.
	if t == anyType {
		return true
	}
	if typ == nil {
		// external call
		typ = f.pkg.types[arg].Type
		if typ == nil {
			return true // probably a type check problem
		}
	}
	// If the type implements fmt.Formatter, we have nothing to check.
	if f.isFormatter(typ) {
		return true
	}
	// If we can use a string, might arg (dynamically) implement the Stringer or Error interface?
	if t&argString != 0 && isConvertibleToString(typ) {
		return true
	}

	typ = typ.Underlying()
	if inProgress[typ] {
		// We're already looking at this type. The call that started it will take care of it.
		return true
	}
	inProgress[typ] = true

	switch typ := typ.(type) {
	case *types.Signature:
		return t&argPointer != 0

	case *types.Map:
		// Recur: map[int]int matches %d.
		return t&argPointer != 0 ||
			(f.matchArgTypeInternal(t, typ.Key(), arg, inProgress) && f.matchArgTypeInternal(t, typ.Elem(), arg, inProgress))

	case *types.Chan:
		return t&argPointer != 0

	case *types.Array:
		// Same as slice.
		if types.Identical(typ.Elem().Underlying(), types.Typ[types.Byte]) && t&argString != 0 {
			return true // %s matches []byte
		}
		// Recur: []int matches %d.
		return t&argPointer != 0 || f.matchArgTypeInternal(t, typ.Elem(), arg, inProgress)

	case *types.Slice:
		// Same as array.
		if types.Identical(typ.Elem().Underlying(), types.Typ[types.Byte]) && t&argString != 0 {
			return true // %s matches []byte
		}
		// Recur: []int matches %d. But watch out for
		//	type T []T
		// If the element is a pointer type (type T[]*T), it's handled fine by the Pointer case below.
		return t&argPointer != 0 || f.matchArgTypeInternal(t, typ.Elem(), arg, inProgress)

	case *types.Pointer:
		// Ugly, but dealing with an edge case: a known pointer to an invalid type,
		// probably something from a failed import.
		if typ.Elem().String() == "invalid type" {
			if *verbose {
				f.Warnf(arg.Pos(), "printf argument %v is pointer to invalid or unknown type", f.gofmt(arg))
			}
			return true // special case
		}
		// If it's actually a pointer with %p, it prints as one.
		if t == argPointer {
			return true
		}
		// If it's pointer to struct, that's equivalent in our analysis to whether we can print the struct.
		if str, ok := typ.Elem().Underlying().(*types.Struct); ok {
			return f.matchStructArgType(t, str, arg, inProgress)
		}
		// Check whether the rest can print pointers.
		return t&argPointer != 0

	case *types.Struct:
		return f.matchStructArgType(t, typ, arg, inProgress)

	case *types.Interface:
		// There's little we can do.
		// Whether any particular verb is valid depends on the argument.
		// The user may have reasonable prior knowledge of the contents of the interface.
		return true

	case *types.Basic:
		switch typ.Kind() {
		case types.UntypedBool,
			types.Bool:
			return t&argBool != 0

		case types.UntypedInt,
			types.Int,
			types.Int8,
			types.Int16,
			types.Int32,
			types.Int64,
			types.Uint,
			types.Uint8,
			types.Uint16,
			types.Uint32,
			types.Uint64,
			types.Uintptr:
			return t&argInt != 0

		case types.UntypedFloat,
			types.Float32,
			types.Float64:
			return t&argFloat != 0

		case types.UntypedComplex,
			types.Complex64,
			types.Complex128:
			return t&argComplex != 0

		case types.UntypedString,
			types.String:
			return t&argString != 0

		case types.UnsafePointer:
			return t&(argPointer|argInt) != 0

		case types.UntypedRune:
			return t&(argInt|argRune) != 0

		case types.UntypedNil:
			return false

		case types.Invalid:
			if *verbose {
				f.Warnf(arg.Pos(), "printf argument %v has invalid or unknown type", f.gofmt(arg))
			}
			return true // Probably a type check problem.
		}
		panic("unreachable")
	}

	return false
}

func isConvertibleToString(typ types.Type) bool {
	if bt, ok := typ.(*types.Basic); ok && bt.Kind() == types.UntypedNil {
		// We explicitly don't want untyped nil, which is
		// convertible to both of the interfaces below, as it
		// would just panic anyway.
		return false
	}
	if types.ConvertibleTo(typ, errorType) {
		return true // via .Error()
	}
	if stringerType != nil && types.ConvertibleTo(typ, stringerType) {
		return true // via .String()
	}
	return false
}

// hasBasicType reports whether x's type is a types.Basic with the given kind.
func (f *File) hasBasicType(x ast.Expr, kind types.BasicKind) bool {
	t := f.pkg.types[x].Type
	if t != nil {
		t = t.Underlying()
	}
	b, ok := t.(*types.Basic)
	return ok && b.Kind() == kind
}

// matchStructArgType reports whether all the elements of the struct match the expected
// type. For instance, with "%d" all the elements must be printable with the "%d" format.
func (f *File) matchStructArgType(t printfArgType, typ *types.Struct, arg ast.Expr, inProgress map[types.Type]bool) bool {
	for i := 0; i < typ.NumFields(); i++ {
		typf := typ.Field(i)
		if !f.matchArgTypeInternal(t, typf.Type(), arg, inProgress) {
			return false
		}
		if t&argString != 0 && !typf.Exported() && isConvertibleToString(typf.Type()) {
			// Issue #17798: unexported Stringer or error cannot be properly fomatted.
			return false
		}
	}
	return true
}

var archSizes = types.SizesFor("gc", build.Default.GOARCH)
