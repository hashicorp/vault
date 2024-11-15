package linodego

/*
Pointer takes a value of any type T and returns a pointer to that value.
Go does not allow directly creating pointers to literals, so Pointer enables
abstraction away the pointer logic.

Example:

		booted := true

		createOpts := linodego.InstanceCreateOptions{
			Booted: &booted,
		}

		can be replaced with

		createOpts := linodego.InstanceCreateOptions{
			Booted: linodego.Pointer(true),
		}
*/

func Pointer[T any](value T) *T {
	return &value
}
