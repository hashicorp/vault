<img src="internal/radius.svg" width="250" align="right">

# radius [![GoDoc](https://godoc.org/layeh.com/radius?status.svg)](https://godoc.org/layeh.com/radius)

a Go (golang) [RADIUS](https://tools.ietf.org/html/rfc2865) client and server implementation

## Installation

    go get -u layeh.com/radius

## Client example

```go
package main

import (
	"context"
	"fmt"

	"layeh.com/radius"
	. "layeh.com/radius/rfc2865"
)

func main() {
	packet := radius.New(radius.CodeAccessRequest, []byte(`secret`))
	UserName_SetString(packet, "tim")
	UserPassword_SetString(packet, "12345")
	response, err := radius.Exchange(context.Background(), packet, "localhost:1812")
	if err != nil {
		panic(err)
	}

	fmt.Println("Code:", response.Code)
}
```

## RADIUS Dictionaries

Included in this package is the command line program `radius-dict-gen`. It can be installed with:

    go get -u layeh.com/radius/cmd/radius-dict-gen

Given a FreeRADIUS dictionary, the program will generate helper functions and types for reading and manipulating RADIUS attributes in a packet. It is recommended that generated code be used for any RADIUS dictionary you would like to consume.

Included in this repository are sub-packages of generated helpers for commonly used RADIUS attributes, including [`rfc2865`](https://godoc.org/layeh.com/radius/rfc2865) and [`rfc2866`](https://godoc.org/layeh.com/radius/rfc2866).

## License

MPL 2.0

## Author

Tim Cooper (<tim.cooper@layeh.com>)
