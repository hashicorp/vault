// See note about this package in kdf.go; these implementations aim to provide
// a PKCS#11 v3.0 CKM_SP800_108_*_KDF compatible interface.
package kdf

import (
	"fmt"
	"hash"
	"io"
)

type feedbackKDF struct {
	// Globally provisioned values
	prf        hash.Hash
	parameters []KBKDFParameter
	keyLengths []int

	// Internal counter for this KDF instance.
	counter uint64
	chain   []byte

	// Internal reference to the present key.
	thisKey int
}

var _ io.Reader = &feedbackKDF{}

// Helper function to add our parameters into the underlying PRF.
func (k *feedbackKDF) addParameters() error {
	for index, _parameter := range k.parameters {
		var data []byte = nil
		switch parameter := _parameter.(type) {
		case CounterVariable:
			data = parameter.Encode(k.counter)
		case DKMLength:
			data = parameter.Encode(k.prf, k.keyLengths)
		case ByteArray:
			data = parameter.Data
		case ChainingVariable:
			data = k.chain
		default:
			// Unreachable assuming assumptions guaranteed by NewFeedback(...) hold.
			return fmt.Errorf("unexpected type for KDF parameter at index %d", index)
		}

		// Write this parameter into the PRF input stream.
		k.prf.Write(data)
	}

	return nil
}

// Read a single key (of pre-specified size) out of the KDF output stream.
func (k *feedbackKDF) Read(outputKey []byte) (int, error) {
	// Assumption: we've not been called too many times (and thus k.thisKey
	// is a valid index into k.keyLengths still).
	if k.thisKey >= len(k.keyLengths) {
		return 0, fmt.Errorf("too many calls to KDF.Read(...)")
	}

	// Assumption: the size of this output key matches the expected size of
	// the present key.
	if len(outputKey)*8 != k.keyLengths[k.thisKey] {
		return 0, fmt.Errorf("size mismatch between expected key length and provided buffer; expected %d bits got %d bits", k.keyLengths[k.thisKey], len(outputKey)*8)
	}

	// Write offset (bytes) in the input parameter
	var prfOutput []byte

	// Division here will take the ceiling, so we have no code outside of this
	// loop that performs the Counter-mode KDF.
	// iterations is equivalent to the N variable,
	iterations := (len(outputKey) + k.prf.Size() - 1) / k.prf.Size()
	for iteration := 0; iteration < iterations; iteration++ {
		// First reset the PRF state; we can't guarantee our caller has done
		// so already.
		k.prf.Reset()

		// Next, add (to the PRF input) all the specified parameters, in the
		// order they were specified.
		err := k.addParameters()
		if err != nil {
			return 0, err
		}

		// Finally, compute the result, placing it in chain, from a nil source.
		// This ensures chain gets the entire PRF output and clears out the
		// previous vaule. We can then append it to prfOutput. This guarantees
		// we have enough space to store the result, even if it might overflow
		// (and thus, extend past) the size of outputKey. Additionally, it stores
		// our chaining value for the next iteration.
		k.chain = k.prf.Sum(nil)
		prfOutput = append(prfOutput, k.chain...)

		// Increment counter.
		k.counter += 1
	}

	// Now that we've collected all key material, increment the key pointer
	// before returning.
	k.thisKey += 1

	if k.thisKey >= len(k.keyLengths) {
		// When this is the last key to be derived from this KDF, clear the
		// value of k.chain so it can get GC'd (hopefully faster than if we
		// had left it). This could contain sensitive key material. In the
		// future, this'd be a good place to use memory erasure from
		// https://github.com/golang/go/issues/21865.
		k.chain = nil
	}

	// Finally, copy from prfOutput into the key. Since:
	//     len(prfOutput) >= len(outputKey)
	// it should hold that the return value of copy is len(outputKey).
	return copy(outputKey, prfOutput), nil
}
