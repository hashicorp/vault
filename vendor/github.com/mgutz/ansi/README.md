# ansi

Package ansi is a small, fast library to create ANSI colored strings and codes.

## Install

This install the color viewer and the package itself

```sh
go get -u github.com/mgutz/ansi/cmd/ansi-mgutz
```

## Example

```go
import "github.com/mgutz/ansi"

// colorize a string, SLOW
msg := ansi.Color("foo", "red+b:white")

// create a closure to avoid recalculating ANSI code compilation
phosphorize := ansi.ColorFunc("green+h:black")
msg = phosphorize("Bring back the 80s!")
msg2 := phospohorize("Look, I'm a CRT!")

// cache escape codes and build strings manually
lime := ansi.ColorCode("green+h:black")
reset := ansi.ColorCode("reset")

fmt.Println(lime, "Bring back the 80s!", reset)
```

Other examples

```go
Color(s, "red")            // red
Color(s, "red+b")          // red bold
Color(s, "red+B")          // red blinking
Color(s, "red+u")          // red underline
Color(s, "red+bh")         // red bold bright
Color(s, "red:white")      // red on white
Color(s, "red+b:white+h")  // red bold on white bright
Color(s, "red+B:white+h")  // red blink on white bright
Color(s, "off")            // turn off ansi codes
```

To view color combinations, from terminal.

```sh
ansi-mgutz
```

## Style format

```go
"foregroundColor+attributes:backgroundColor+attributes"
```

Colors

* black
* red
* green
* yellow
* blue
* magenta
* cyan
* white
* 0...255 (256 colors)

Attributes

*   b = bold foreground
*   B = Blink foreground
*   u = underline foreground
*   i = inverse
*   h = high intensity (bright) foreground, background

    does not work with 256 colors

## Constants

* ansi.Reset
* ansi.DefaultBG
* ansi.DefaultFG
* ansi.Black
* ansi.Red
* ansi.Green
* ansi.Yellow
* ansi.Blue
* ansi.Magenta
* ansi.Cyan
* ansi.White
* ansi.LightBlack
* ansi.LightRed
* ansi.LightGreen
* ansi.LightYellow
* ansi.LightBlue
* ansi.LightMagenta
* ansi.LightCyan
* ansi.LightWhite


## References

Wikipedia ANSI escape codes [Colors](http://en.wikipedia.org/wiki/ANSI_escape_code#Colors)

General [tips and formatting](http://misc.flogisoft.com/bash/tip_colors_and_formatting)

What about support on Windows? Use [colorable by mattn](https://github.com/mattn/go-colorable).
Ansi and colorable are used by [logxi](https://github.com/mgutz/logxi) to support logging in
color on Windows.

## MIT License

Copyright (c) 2013 Mario Gutierrez mario@mgutz.com

See the file LICENSE for copying permission.

