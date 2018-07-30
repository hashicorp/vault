package parseutil

import (
	"net/http"
	"net/textproto"
)

/*
	The reason we wrap an http.Header here is - under the hood,
	an http.Header is a map[string][]string and it does some magic
	with casing that can make it difficult to correctly iterate
	all header values.

	For example:
	h := http.Header{}
	h.Add("hello", "world")
	h.Add("hello", "monde")

	You'd expect that the header now have "hello" for a key, and
	both "world" and "monde" for a value. But it doesn't. "Hello"
	is now capitalized in the map, but "world" and "monde" aren't.

	Later, when you do this:
	h.Get("hello")

	You'll receive only "world", silently missing "monde".

	If you try to solve that by doing this:
	h["hello"]

	You'll receive nothing back, because remember, it's now in the
	map as "Hello".

	To avoid bugs like this, we provide one more method. GetAll,
	which returns all the values for a key.
*/
func NewHeader() Header {
	return Header{http.Header{}}
}

type Header struct {
	http.Header
}

func (h Header) GetAll(key string) []string {
	return h.Header[textproto.CanonicalMIMEHeaderKey(key)]
}
