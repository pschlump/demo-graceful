package HashStrings512

//
// Test of sha512 hashing
//

import (
	"fmt"
	"testing"
)

func Test_Hash512(t *testing.T) {

	hh := Sha512("angryMonkey")
	// expect := "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7s u2A+gf7Q=="
	expect := "6441e1581eb9814973755c2d0d002b132c7e2952f3a7f69369168f941cd8448163eaf8c576a11bd10e41f3354a099d2f29b64f664949cf415deecbb603e81fed"
	if hh != expect {
		t.Errorf("Unexpecd value for hash, got [%s] expected [%s]\n", hh, expect)
	}

	hh = HashStrings512("angry", "Monkey")
	if hh != expect {
		t.Errorf("Unexpecd value for hash, got [%s] expected [%s]\n", hh, expect)
	}

	bb := HashByte512([]byte("angryMonkey"))
	hh = fmt.Sprintf("%x", bb)
	if hh != expect {
		t.Errorf("Unexpecd value for hash, got [%s] expected [%s]\n", bb, expect)
	}

	bb = HashBytes512([]byte("angry"), []byte("Monkey"))
	hh = fmt.Sprintf("%x", bb)
	if hh != expect {
		t.Errorf("Unexpecd value for hash, got [%s] expected [%s]\n", bb, expect)
	}
}

/* vim: set noai ts=4 sw=4: */
