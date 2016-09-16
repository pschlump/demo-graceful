package HashStrings512

import (
	"crypto/sha512"
	"fmt"
)

// Single item interface
func Sha512(s string) (rv string) {
	rv = HashStrings512(s)
	return
}

// hash a set of strings and return in hex-strings form
func HashStrings512(a ...string) string {
	h := sha512.New()
	for _, z := range a {
		h.Write([]byte(z))
	}
	return fmt.Sprintf("%x", (h.Sum(nil)))
}

func HashByte512(a []byte) (rv []byte) {
	rv = HashBytes512(a)
	return
}

// hash a set of []byte and return in byte form
func HashBytes512(a ...[]byte) []byte {
	h := sha512.New()
	for _, z := range a {
		h.Write(z)
	}
	return h.Sum(nil)
}
