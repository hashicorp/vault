# bexpr - Boolean Expression Evaluator [![GoDoc](https://godoc.org/github.com/hashicorp/go-bexpr?status.svg)](https://godoc.org/github.com/hashicorp/go-bexpr) [![CircleCI](https://circleci.com/gh/hashicorp/go-bexpr.svg?style=svg)](https://circleci.com/gh/hashicorp/go-bexpr)

`bexpr` is a Go (golang) library to provide generic boolean expression
evaluation and filtering for Go data structures and maps. Under the hood,
`bexpr` uses
[`pointerstructure`](https://github.com/mitchellh/pointerstructure), meaning
that any path within a map or structure that can be expressed via that library
can be used with `bexpr`. This also means that you can use the custom `bexpr`
dotted syntax (kept mainly for backwards compatibility) to select values in
expressions, or, by enclosing the selectors in quotes, you can use [JSON
Pointer](https://tools.ietf.org/html/rfc6901) syntax to select values in
expressions.

## Usage (Reflection)

This example program is available in [examples/simple](examples/simple)

```go
package main

import (
   "fmt"
   "github.com/hashicorp/go-bexpr"
)

type Example struct {
   X int

   // Can rename a field with the struct tag
   Y string `bexpr:"y"`
   Z bool `bexpr:"foo"`

   // Tag with "-" to prevent allowing this field from being used
   Hidden string `bexpr:"-"`

   // Unexported fields are not available for evaluation
   unexported string
}

func main() {
   value := map[string]Example{
      "foo": Example{X: 5, Y: "foo", Z: true, Hidden: "yes", unexported: "no"},
      "bar": Example{X: 42, Y: "bar", Z: false, Hidden: "no", unexported: "yes"},
   }

   expressions := []string{
		"foo.X == 5",
		"bar.y == bar",
		"foo.baz == true",

		// will error in evaluator creation
		"bar.Hidden != yes",

		// will error in evaluator creation
		"foo.unexported == no",
	}

   for _, expression := range expressions {
      eval, err := bexpr.CreateEvaluator(expression)

      if err != nil {
         fmt.Printf("Failed to create evaluator for expression %q: %v\n", expression, err)
         continue
      }

      result, err := eval.Evaluate(value)
      if err != nil {
         fmt.Printf("Failed to run evaluation of expression %q: %v\n", expression, err)
         continue
      }

      fmt.Printf("Result of expression %q evaluation: %t\n", expression, result)
   }
}
```

This will output:

```
Result of expression "foo.X == 5" evaluation: true
Result of expression "bar.y == bar" evaluation: true
Result of expression "foo.baz == true" evaluation: true
Failed to run evaluation of expression "bar.Hidden != yes": error finding value in datum: /bar/Hidden at part 1: struct field "Hidden" is ignored and cannot be used
Failed to run evaluation of expression "foo.unexported == no": error finding value in datum: /foo/unexported at part 1: couldn't find struct field with name "unexported"
```

## Testing

The [Makefile](Makefile) contains 3 main targets to aid with testing:

1. `make test` - runs the standard test suite
2. `make coverage` - runs the test suite gathering coverage information
3. `make bench` - this will run benchmarks. You can use the [`benchcmp`](https://godoc.org/golang.org/x/tools/cmd/benchcmp) tool to compare
   subsequent runs of the tool to compare performance. There are a few arguments you can
   provide to the make invocation to alter the behavior a bit
   * `BENCHFULL=1` - This will enable running all the benchmarks. Some could be fairly redundant but
     could be useful when modifying specific sections of the code.
   * `BENCHTIME=5s` - By default the -benchtime paramater used for the `go test` invocation is `2s`.
     `1s` seemed like too little to get results consistent enough for comparison between two runs.
     For the highest degree of confidence that performance has remained steady increase this value
     even further. The time it takes to run the bench testing suite grows linearly with this value.
   * `BENCHTESTS=BenchmarkEvaluate` - This is used to run a particular benchmark including all of its
     sub-benchmarks. This is just an example and "BenchmarkEvaluate" can be replaced with any
     benchmark functions name.
