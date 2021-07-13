package base62

import (
	"fmt"
)

/*MIT License

Copyright (c) 2017 Denis Subbotin
Copyright (c) 2017 Nika Jones
Copyright (c) 2017 Philip Schlump
Copyright (c) 2020 Peter 'ribasushi' Rabbitson

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

Credit to github.com/mr-tron/base58 for the implementation
*/

var lookup [128]int8

func init() {
	for i, c := range charset {
		lookup[c]= int8(i)
	}
}


// EncodeToString encodes bytes to base62.  This does *not* scale linearly with input as base64, so use caution
// when using on large inputs.
func EncodeToString(str []byte) string {
	if str == nil {
		return ""
	}


	size := len(str)

	zcount := 0
	for zcount < size && str[zcount] == 0 {
		zcount++
	}

	// It is crucial to make this as short as possible, especially for
	// the usual case of bitcoin addrs
	size = zcount +
		// This is an integer simplification of
		// ceil(log(256)/log(62))
		(size-zcount)*262/195 + 1


	out := make([]byte, size)

	var i, high int
	var carry uint32

	high = size - 1
	for _, b := range str {
		i = size - 1
		for carry = uint32(b); i > high || carry != 0; i-- {
			carry = carry + 256*uint32(out[i])
			out[i] = byte(carry % 62)
			carry /= 62
		}
		high = i
	}

	// Determine the additional "zero-gap" in the buffer (aside from zcount)
	for i = zcount; i < size && out[i] == 0; i++ {
	}

	// Now encode the values with actual alphabet in-place
	val := out[i-zcount:]
	size = len(val)
	for i = 0; i < size; i++ {
		out[i] = charset[val[i]]
	}

	return string(out[:size])
}

// Decode decodes a base62 string into bytes. This does *not* scale linearly with input as base64, so use caution
// when using on large inputs.
func DecodeString(str string) ([]byte, error) {
	if str=="" {
		return nil, nil
	}

	zero := uint8(lookup[0])
	b62sz := len(str)

	var zcount int
	for i := 0; i < b62sz && str[i] == zero; i++ {
		zcount++
	}

	var t, c uint64

	// the 32bit algo stretches the result up to 2 times
	binu := make([]byte, (b62sz*195/262 + 1)<<1)
	outi := make([]uint32, (b62sz+3)/4)

	for _, r := range str {
		if r > 127 {
			return nil, fmt.Errorf("high-bit set on invalid digit")
		}
		if lookup[r] == -1 {
			return nil, fmt.Errorf("invalid base62 digit (%q)", r)
		}

		c = uint64(lookup[r])

		for j := len(outi) - 1; j >= 0; j-- {
			t = uint64(outi[j])*62 + c
			c = t >> 32
			outi[j] = uint32(t & 0xffffffff)
		}
	}

	// initial mask depends on b62sz, on further loops it always starts at 24 bits
	mask := uint(b62sz%4) * 8
	if mask == 0 {
		mask = 32
	}
	mask -= 8

	outLen := 0
	for j := 0; j < len(outi); j++ {
		for mask < 32 { // loop relies on uint overflow
			binu[outLen] = byte(outi[j] >> mask)
			mask -= 8
			outLen++
		}
		mask = 24
	}

	// find the most significant byte post-decode, if any
	for msb := zcount; msb < len(binu); msb++ {
		if binu[msb] > 0 {
			return binu[msb-zcount : outLen], nil
		}
	}

	// it's all zeroes
	return binu[:outLen], nil
}