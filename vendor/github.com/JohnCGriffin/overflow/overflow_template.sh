#!/bin/sh

exec > overflow_impl.go

echo "package overflow

// This is generated code, created by overflow_template.sh executed
// by \"go generate\"

"


for SIZE in 8 16 32 64
do
echo "

// Add${SIZE} performs + operation on two int${SIZE} operands
// returning a result and status
func Add${SIZE}(a, b int${SIZE}) (int${SIZE}, bool) {
        c := a + b
        if (c > a) == (b > 0) {
                return c, true
        }
        return c, false
}

// Add${SIZE}p is the unchecked panicing version of Add${SIZE}
func Add${SIZE}p(a, b int${SIZE}) int${SIZE} {
        r, ok := Add${SIZE}(a, b)
        if !ok {
                panic(\"addition overflow\")
        }
        return r
}


// Sub${SIZE} performs - operation on two int${SIZE} operands
// returning a result and status
func Sub${SIZE}(a, b int${SIZE}) (int${SIZE}, bool) {
        c := a - b
        if (c < a) == (b > 0) {
                return c, true
        }
        return c, false
}

// Sub${SIZE}p is the unchecked panicing version of Sub${SIZE}
func Sub${SIZE}p(a, b int${SIZE}) int${SIZE} {
        r, ok := Sub${SIZE}(a, b)
        if !ok {
                panic(\"subtraction overflow\")
        }
        return r
}


// Mul${SIZE} performs * operation on two int${SIZE} operands
// returning a result and status
func Mul${SIZE}(a, b int${SIZE}) (int${SIZE}, bool) {
        if a == 0 || b == 0 {
                return 0, true
        }
        c := a * b
        if (c < 0) == ((a < 0) != (b < 0)) {
                if c/b == a {
                        return c, true
                }
        }
        return c, false
}

// Mul${SIZE}p is the unchecked panicing version of Mul${SIZE}
func Mul${SIZE}p(a, b int${SIZE}) int${SIZE} {
        r, ok := Mul${SIZE}(a, b)
        if !ok {
                panic(\"multiplication overflow\")
        }
        return r
}



// Div${SIZE} performs / operation on two int${SIZE} operands
// returning a result and status
func Div${SIZE}(a, b int${SIZE}) (int${SIZE}, bool) {
        q, _, ok := Quotient${SIZE}(a, b)
        return q, ok
}

// Div${SIZE}p is the unchecked panicing version of Div${SIZE}
func Div${SIZE}p(a, b int${SIZE}) int${SIZE} {
        r, ok := Div${SIZE}(a, b)
        if !ok {
                panic(\"division failure\")
        }
        return r
}

// Quotient${SIZE} performs + operation on two int${SIZE} operands
// returning a quotient, a remainder and status
func Quotient${SIZE}(a, b int${SIZE}) (int${SIZE}, int${SIZE}, bool) {
        if b == 0 {
                return 0, 0, false
        }
        c := a / b
        status := (c < 0) == ((a < 0) != (b < 0))
        return c, a % b, status
}
"
done
