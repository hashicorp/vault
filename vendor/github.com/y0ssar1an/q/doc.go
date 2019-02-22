/*
Package q provides quick and dirty debugging output for tired programmers.

q.Q() is a fast way to pretty-print variables. It's easier than typing
fmt.Printf("%#v", whatever). The output will be colorized and nicely formatted.
The output goes to $TMPDIR/q, away from the noise of stdout.

q exports a single Q() function. This is how you use it:
    import "github.com/y0ssar1an/q"
    ...
    q.Q(a, b, c)
*/
package q
