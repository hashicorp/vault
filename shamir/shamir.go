package shamir

import (
	"crypto/rand"
	"crypto/subtle"
	"fmt"
	mathrand "math/rand"
	"time"

	"github.com/hashicorp/errwrap"
)

const (
	// ShareOverhead is the byte size overhead of each share
	// when using Split on a secret. This is caused by appending
	// a one byte tag to the share.
	ShareOverhead = 1
)

// polynomial represents a polynomial of arbitrary degree
type polynomial struct {
	coefficients []uint8
}

// makePolynomial constructs a random polynomial of the given
// degree but with the provided intercept value.
func makePolynomial(intercept, degree uint8) (polynomial, error) {
	// Create a wrapper
	p := polynomial{
		coefficients: make([]byte, degree+1),
	}

	// Ensure the intercept is set
	p.coefficients[0] = intercept

	// Assign random co-efficients to the polynomial
	if _, err := rand.Read(p.coefficients[1:]); err != nil {
		return p, err
	}

	return p, nil
}

// evaluate returns the value of the polynomial for the given x
func (p *polynomial) evaluate(x uint8) uint8 {
	// Special case the origin
	if x == 0 {
		return p.coefficients[0]
	}

	// Compute the polynomial value using Horner's method.
	degree := len(p.coefficients) - 1
	out := p.coefficients[degree]
	for i := degree - 1; i >= 0; i-- {
		coeff := p.coefficients[i]
		out = add(mult(out, x), coeff)
	}
	return out
}

// interpolatePolynomial takes N sample points and returns
// the value at a given x using a lagrange interpolation.
func interpolatePolynomial(x_samples, y_samples []uint8, x uint8) uint8 {
	limit := len(x_samples)
	var result, basis uint8
	for i := 0; i < limit; i++ {
		basis = 1
		for j := 0; j < limit; j++ {
			if i == j {
				continue
			}
			num := add(x, x_samples[j])
			denom := add(x_samples[i], x_samples[j])
			term := div(num, denom)
			basis = mult(basis, term)
		}
		group := mult(y_samples[i], basis)
		result = add(result, group)
	}
	return result
}

// div divides two numbers in GF(2^8)
func div(a, b uint8) uint8 {
	if b == 0 {
		// leaks some timing information but we don't care anyways as this
		// should never happen, hence the panic
		panic("divide by zero")
	}

	var goodVal, zero uint8
	log_a := logTable[a]
	log_b := logTable[b]
	diff := (int(log_a) - int(log_b)) % 255
	if diff < 0 {
		diff += 255
	}

	ret := expTable[diff]

	// Ensure we return zero if a is zero but aren't subject to timing attacks
	goodVal = ret

	if subtle.ConstantTimeByteEq(a, 0) == 1 {
		ret = zero
	} else {
		ret = goodVal
	}

	return ret
}

// mult multiplies two numbers in GF(2^8)
func mult(a, b uint8) (out uint8) {
	var goodVal, zero uint8
	log_a := logTable[a]
	log_b := logTable[b]
	sum := (int(log_a) + int(log_b)) % 255

	ret := expTable[sum]

	// Ensure we return zero if either a or b are zero but aren't subject to
	// timing attacks
	goodVal = ret

	if subtle.ConstantTimeByteEq(a, 0) == 1 {
		ret = zero
	} else {
		ret = goodVal
	}

	if subtle.ConstantTimeByteEq(b, 0) == 1 {
		ret = zero
	} else {
		// This operation does not do anything logically useful. It
		// only ensures a constant number of assignments to thwart
		// timing attacks.
		goodVal = zero
	}

	return ret
}

// add combines two numbers in GF(2^8)
// This can also be used for subtraction since it is symmetric.
func add(a, b uint8) uint8 {
	return a ^ b
}

// Split takes an arbitrarily long secret and generates a `parts`
// number of shares, `threshold` of which are required to reconstruct
// the secret. The parts and threshold must be at least 2, and less
// than 256. The returned shares are each one byte longer than the secret
// as they attach a tag used to reconstruct the secret.
func Split(secret []byte, parts, threshold int) ([][]byte, error) {
	// Sanity check the input
	if parts < threshold {
		return nil, fmt.Errorf("parts cannot be less than threshold")
	}
	if parts > 255 {
		return nil, fmt.Errorf("parts cannot exceed 255")
	}
	if threshold < 2 {
		return nil, fmt.Errorf("threshold must be at least 2")
	}
	if threshold > 255 {
		return nil, fmt.Errorf("threshold cannot exceed 255")
	}
	if len(secret) == 0 {
		return nil, fmt.Errorf("cannot split an empty secret")
	}

	// Generate random list of x coordinates
	mathrand.Seed(time.Now().UnixNano())
	xCoordinates := mathrand.Perm(255)

	// Allocate the output array, initialize the final byte
	// of the output with the offset. The representation of each
	// output is {y1, y2, .., yN, x}.
	out := make([][]byte, parts)
	for idx := range out {
		out[idx] = make([]byte, len(secret)+1)
		out[idx][len(secret)] = uint8(xCoordinates[idx]) + 1
	}

	// Construct a random polynomial for each byte of the secret.
	// Because we are using a field of size 256, we can only represent
	// a single byte as the intercept of the polynomial, so we must
	// use a new polynomial for each byte.
	for idx, val := range secret {
		p, err := makePolynomial(val, uint8(threshold-1))
		if err != nil {
			return nil, errwrap.Wrapf("failed to generate polynomial: {{err}}", err)
		}

		// Generate a `parts` number of (x,y) pairs
		// We cheat by encoding the x value once as the final index,
		// so that it only needs to be stored once.
		for i := 0; i < parts; i++ {
			x := uint8(xCoordinates[i]) + 1
			y := p.evaluate(x)
			out[i][idx] = y
		}
	}

	// Return the encoded secrets
	return out, nil
}

// Combine is used to reverse a Split and reconstruct a secret
// once a `threshold` number of parts are available.
func Combine(parts [][]byte) ([]byte, error) {
	// Verify enough parts provided
	if len(parts) < 2 {
		return nil, fmt.Errorf("less than two parts cannot be used to reconstruct the secret")
	}

	// Verify the parts are all the same length
	firstPartLen := len(parts[0])
	if firstPartLen < 2 {
		return nil, fmt.Errorf("parts must be at least two bytes")
	}
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) != firstPartLen {
			return nil, fmt.Errorf("all parts must be the same length")
		}
	}

	// Create a buffer to store the reconstructed secret
	secret := make([]byte, firstPartLen-1)

	// Buffer to store the samples
	x_samples := make([]uint8, len(parts))
	y_samples := make([]uint8, len(parts))

	// Set the x value for each sample and ensure no x_sample values are the same,
	// otherwise div() can be unhappy
	checkMap := map[byte]bool{}
	for i, part := range parts {
		samp := part[firstPartLen-1]
		if exists := checkMap[samp]; exists {
			return nil, fmt.Errorf("duplicate part detected")
		}
		checkMap[samp] = true
		x_samples[i] = samp
	}

	// Reconstruct each byte
	for idx := range secret {
		// Set the y value for each sample
		for i, part := range parts {
			y_samples[i] = part[idx]
		}

		// Interpolate the polynomial and compute the value at 0
		val := interpolatePolynomial(x_samples, y_samples, 0)

		// Evaluate the 0th value to get the intercept
		secret[idx] = val
	}
	return secret, nil
}
