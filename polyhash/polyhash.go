package polyhash

import (
	"crypto/rand"
	"crypto/sha256"

	"fmt"
)

const SHARE_LENGTH = 32
const SALT_LENGTH = 16

func check_error(e error) {
	if e != nil {
		panic(e)
	}
}

// xorBytes returns the xor'ed value of it's inputs
func xorBytes(a, b []byte) []byte {
	n := len(a)
	dst := make([]byte, n)
	for i := 0; i < n; i++ {
		dst[i] = a[i] ^ b[i]
	}
	return dst
}

// computeHash helps to compute the salted hash from the password
func computeHash(salt []byte, password string) []byte {
	hash_func := sha256.New()
	hash_func.Write(salt)
	hash_func.Write([]byte(password))
	return hash_func.Sum(nil)
}

// StoreShareInformation takes in list of passwords, and a list of shares and
// returns a list of polyhashes of the form "share number,share xor hash,salt"
func StoreShareInformation(passwords []string, shares [][]byte) []string {

	var shareno int = 1
	salt := make([]byte, SALT_LENGTH)
	var polyhashentry string
	var polyhashdb []string
	var err error

	for i := 0; i < len(shares); i++ {
		// Generate a random salt
		_, err = rand.Read(salt)
		check_error(err)

		// Compute the share_xor_hashes
		hash := computeHash(salt, passwords[i])
		share_xor_hash := xorBytes(shares[i][:len(shares[i])-1], hash)

		// FIXME probably concatenating this tring is not the best idea
		polyhashentry = fmt.Sprintf("%02x,%064x,%032x\n", shareno,
			share_xor_hash, salt)

		polyhashdb = append(polyhashdb, polyhashentry)
		shareno += 1
	}

	return polyhashdb
}

// RecoverShareFromPolyhash recovers the share value from polyhashentry
// using the share number and password that are passed in
func RecoverShareFromPolyhash(share_num int, polyhashdb []string, password string) []byte {

	share_xor_hash := make([]byte, SHARE_LENGTH)
	salt := make([]byte, SALT_LENGTH)
	var number int

	for _, entry := range polyhashdb {
		fmt.Sscanf(entry, "%x,%x,%x", &number, &share_xor_hash, &salt)
		if number == share_num {
			break
		}
	}

	hash := computeHash(salt, password)
	share := xorBytes(share_xor_hash, hash)
	shareno := byte(share_num)
	share = append(share, shareno)
	return share
}
