package rfc3961

/*
Implementation of the n-fold algorithm as defined in RFC 3961.

n-fold is an algorithm that takes m input bits and "stretches" them
to form n output bits with equal contribution from each input bit to
the output, as described in [Blumenthal96]:

We first define a primitive called n-folding, which takes a
variable-length input block and produces a fixed-length output
sequence.  The intent is to give each input bit approximately
equal weight in determining the value of each output bit.  Note
that whenever we need to treat a string of octets as a number, the
assumed representation is Big-Endian -- Most Significant Byte
first.

To n-fold a number X, replicate the input value to a length that
is the least common multiple of n and the length of X.  Before
each repetition, the input is rotated to the right by 13 bit
positions.  The successive n-bit chunks are added together using
1's-complement addition (that is, with end-around carry) to yield
a n-bit result....
*/

/* Credits
This golang implementation of nfold used the following project for help with implementation detail.
Although their source is in java it was helpful as a reference implementation of the RFC.
You can find the source code of their open source project along with license information below.
We acknowledge and are grateful to these developers for their contributions to open source

Project: Apache Directory (http://http://directory.apache.org/)
https://svn.apache.org/repos/asf/directory/apacheds/tags/1.5.1/kerberos-shared/src/main/java/org/apache/directory/server/kerberos/shared/crypto/encryption/NFold.java
License: http://www.apache.org/licenses/LICENSE-2.0
*/

// Nfold expands the key to ensure it is not smaller than one cipher block.
// Defined in RFC 3961.
//
// m input bytes that will be "stretched" to the least common multiple of n bits and the bit length of m.
func Nfold(m []byte, n int) []byte {
	k := len(m) * 8

	//Get the lowest common multiple of the two bit sizes
	lcm := lcm(n, k)
	relicate := lcm / k
	var sumBytes []byte

	for i := 0; i < relicate; i++ {
		rotation := 13 * i
		sumBytes = append(sumBytes, rotateRight(m, rotation)...)
	}

	nfold := make([]byte, n/8)
	sum := make([]byte, n/8)
	for i := 0; i < lcm/n; i++ {
		for j := 0; j < n/8; j++ {
			sum[j] = sumBytes[j+(i*len(sum))]
		}
		nfold = onesComplementAddition(nfold, sum)
	}
	return nfold
}

func onesComplementAddition(n1, n2 []byte) []byte {
	numBits := len(n1) * 8
	out := make([]byte, numBits/8)
	carry := 0
	for i := numBits - 1; i > -1; i-- {
		n1b := getBit(&n1, i)
		n2b := getBit(&n2, i)
		s := n1b + n2b + carry

		if s == 0 || s == 1 {
			setBit(&out, i, s)
			carry = 0
		} else if s == 2 {
			carry = 1
		} else if s == 3 {
			setBit(&out, i, 1)
			carry = 1
		}
	}
	if carry == 1 {
		carryArray := make([]byte, len(n1))
		carryArray[len(carryArray)-1] = 1
		out = onesComplementAddition(out, carryArray)
	}
	return out
}

func rotateRight(b []byte, step int) []byte {
	out := make([]byte, len(b))
	bitLen := len(b) * 8
	for i := 0; i < bitLen; i++ {
		v := getBit(&b, i)
		setBit(&out, (i+step)%bitLen, v)
	}
	return out
}

func lcm(x, y int) int {
	return (x * y) / gcd(x, y)
}

func gcd(x, y int) int {
	for y != 0 {
		x, y = y, x%y
	}
	return x
}

func getBit(b *[]byte, p int) int {
	pByte := p / 8
	pBit := uint(p % 8)
	vByte := (*b)[pByte]
	vInt := int(vByte >> (8 - (pBit + 1)) & 0x0001)
	return vInt
}

func setBit(b *[]byte, p, v int) {
	pByte := p / 8
	pBit := uint(p % 8)
	oldByte := (*b)[pByte]
	var newByte byte
	newByte = byte(v<<(8-(pBit+1))) | oldByte
	(*b)[pByte] = newByte
}
