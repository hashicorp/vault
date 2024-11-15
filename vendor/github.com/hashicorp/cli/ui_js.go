package cli

import (
	"syscall/js"
)

func (u *BasicUi) ask(query string, secret bool) (string, error) {
	line := js.Global().Call("prompt", query).String()
	return line, nil
}
