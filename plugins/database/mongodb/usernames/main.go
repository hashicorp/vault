package main

import (
	"fmt"
	"strings"
	"text/template"
)

// type Inventory struct {
// 	Material string
// 	Count    uint
//
// 	Value   interface{}
// 	Mapping map[string]interface{}
// 	Slice   []string
// }
//
// type SubType struct {
// 	Foo string
// }
//
// func getCallerInfo() (funcName, file string, line int, err error) {
// 	const callDepth = 2 // user code calls q.Q() which calls std.log().
// 	pc, file, line, ok := runtime.Caller(callDepth)
// 	if !ok {
// 		// This error is not exported. It is only used internally in the q
// 		// package. The error message isn't even used by the caller. So, I've
// 		// suppressed the goerr113 linter here, which catches nonidiomatic
// 		// error handling post Go 1.13 errors.
// 		return "", "", 0, errors.New("failed to get info about the function calling q.Q") // nolint: goerr113
// 	}
//
// 	funcName = runtime.FuncForPC(pc).Name()
// 	return funcName, file, line, nil
// }
//
// func argNames(filename string, line int) ([]string, error) {
// 	fset := token.NewFileSet()
// 	f, err := parser.ParseFile(fset, filename, nil, 0)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to parse %q: %w", filename, err)
// 	}
//
// 	var names []string
// 	ast.Inspect(f, func(n ast.Node) bool { // nolint: unparam
// 		call, is := n.(*ast.CallExpr)
// 		if !is {
// 			// The node is not a function call.
// 			return true // visit next node
// 		}
//
// 		if fset.Position(call.End()).Line != line {
// 			// The node is a function call, but it's on the wrong line.
// 			return true
// 		}
//
// 		// if !isQCall(call) {
// 		// 	// The node is a function call on correct line, but it's not a Q()
// 		// 	// function.
// 		// 	return true
// 		// }
//
// 		for _, arg := range call.Args {
// 			names = append(names, argName(arg))
// 		}
// 		return true
// 	})
//
// 	return names, nil
// }
//
// func argName(arg ast.Expr) string {
// 	name := ""
//
// 	switch a := arg.(type) {
// 	case *ast.Ident:
// 		switch {
// 		case a.Obj == nil:
// 			name = a.Name
// 		case a.Obj.Kind == ast.Var, a.Obj.Kind == ast.Con:
// 			name = a.Obj.Name
// 		}
// 	case *ast.BinaryExpr,
// 		*ast.CallExpr,
// 		*ast.IndexExpr,
// 		*ast.KeyValueExpr,
// 		*ast.ParenExpr,
// 		*ast.SelectorExpr,
// 		*ast.SliceExpr,
// 		*ast.TypeAssertExpr,
// 		*ast.UnaryExpr:
// 		name = exprToString(arg)
// 	}
//
// 	return name
// }
//
// func exprToString(arg ast.Expr) string {
// 	var buf strings.Builder
// 	fset := token.NewFileSet()
// 	if err := printer.Fprint(&buf, fset, arg); err != nil {
// 		return ""
// 	}
//
// 	// CallExpr will be multi-line and indented with tabs. replace tabs with
// 	// spaces so we can better control formatting during output().
// 	return strings.Replace(buf.String(), "\t", "    ", -1)
// }
//
// func FormatUsername(tmpl string, data ...interface{}) (username string, err error) {
// 	funcName, file, line, err := getCallerInfo()
// 	if err != nil {
// 		return "", fmt.Errorf("unable to get caller info: %w", err)
// 	}
// 	fmt.Printf("Func: %s\nFile: %s\nLine: %d\n", funcName, file, line)
//
// 	args, err := argNames(file, line)
// 	if err != nil {
// 		return "", fmt.Errorf("unable to get argument names: %w", err)
// 	}
// 	fmt.Printf("Args: %s\n", strings.Join(args, "; "))
//
// 	return "", nil
// }

func main() {
	sandbox()

	// userConfiguredTemplate := ""
	//
	// username, err := FormatUsername(userConfiguredTemplate)
	// if err != nil {
	// 	fmt.Printf("ERROR: %s\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Printf("Username: %s\n", username)
}

func FormatUsername(rawTemplate string, data interface{}) (username string, err error) {
	tmpl, err := template.
		New("username").
		Option("missingkey=error").
		Parse(rawTemplate)
	if err != nil {
		return "", fmt.Errorf("unable to parse template string: %w", err)
	}

	str := &strings.Builder{}
	err = tmpl.Execute(str, data)
	if err != nil {
		return "", fmt.Errorf("failed to process template: %w", err)
	}

	return str.String(), nil
}

type EndpointData struct {
	DisplayName string
	RoleName    string
}

type DBSpecificData struct {
	EndpointData
	Foo string
}

func sandbox() {
	data := DBSpecificData{
		EndpointData: EndpointData{
			DisplayName: "displayname",
			RoleName:    "rolename",
		},
		Foo: "foo",
	}

	username, err := generateUsername("{{.DisplayName | truncate 5}}-{{.RoleName | truncate 5}}-{{.Foo | truncate 5}}-{{now_seconds}}-{{rand 20}}", data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Username: %s\n", username)
}

func generateUsername(userSpecifiedTemplate string, dataSpecificToEngine interface{}) (username string, err error) {
	producer, err := NewUsernameProducer(
		Template(userSpecifiedTemplate),
		MaxLength(30),
	)
	if err != nil {
		return "", err
	}
	username, err = producer.GenerateUsername(dataSpecificToEngine)
	if err != nil {
		return "", err
	}
	return username, nil
}
