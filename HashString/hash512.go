package HashStrings512

import (
	"crypto/sha512"
	"fmt"
)

//

// Sha512 takes a string and returns the sha512 hash of that string in hex
func Sha512(s string) (rv string) {
	rv = HashStrings512(s)
	return
}

// HashString512 hashs a set of strings and return in hex-strings form the sha512 hash
func HashStrings512(a ...string) string {
	h := sha512.New()
	for _, z := range a {
		h.Write([]byte(z))
	}
	return fmt.Sprintf("%x", (h.Sum(nil)))
}

// HashBytes512 takes a []byte and returns the sha512 for the array
func HashByte512(a []byte) (rv []byte) {
	rv = HashBytes512(a)
	return
}

// HashBytess512 takes a set of []byte and returns the sha512 for the concatenated set
func HashBytes512(a ...[]byte) []byte {
	h := sha512.New()
	for _, z := range a {
		h.Write(z)
	}
	return h.Sum(nil)
}

/* vim: set noai ts=4 sw=4: */
